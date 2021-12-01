package sqlite

import (
	"context"
	"errors"
	"log"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/google/uuid"
)

type PortfolioService struct {
	db *DB
}

func NewPortfolioService(db *DB) *PortfolioService {
	return &PortfolioService{db}
}

func (s *PortfolioService) FindPortfolioByID(ctx context.Context, id uuid.UUID) (*app.Portfolio, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	portfolio, err := findPortfolioByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return portfolio, tx.Commit()
}

func findPortfolioByID(ctx context.Context, tx *Tx, id uuid.UUID) (*app.Portfolio, error) {
	portfolio := &app.Portfolio{}
	var userId string
	err := tx.QueryRowContext(ctx, `
		SELECT
			portfolio_id,
			title,
			owner_id,
			(SELECT user_id FROM providers WHERE provider_id = owner_id) AS user_id
		FROM portfolios
		WHERE portfolio_id = ?
	`, id).Scan(
		&portfolio.ID,
		&portfolio.Title,
		&portfolio.ProviderID,
		&userId,
	)
	if err != nil {
		return nil, err
	}

	// Get photos
	rows, err := tx.QueryContext(ctx, `
		SELECT
			photo_url
		FROM photos
		WHERE uploaded_by = ?
		`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var photo string
		if err := rows.Scan(
			&photo,
		); err != nil {
			return nil, err
		}
		portfolio.Photos = append(portfolio.Photos, photo)
	}

	return portfolio, nil
}

func (s *PortfolioService) CreatePortfolio(ctx context.Context, portfolio *model.Portfolio) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	/*
		provider, err := getProviderByCriteria(ctx, tx, "user_id", portfolio.UserID)
		if err != nil {
			log.Println("getProviderByCriteria failed:", err)
			return err
		}

		portfolio.OwnerID = provider.ID
	*/
	// Create dial and attach associated owner user.
	if err := createPortfolio(ctx, tx, portfolio); err != nil {
		return err
	}
	return tx.Commit()
}

func createPortfolio(ctx context.Context, tx *Tx, portfolio *model.Portfolio) error {
	portfolio.ID = uuid.New()

	// Create dial and attach associated owner user.
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO portfolios (
			portfolio_id,
			owner_id,
			title,
			booking_id,
			service_id
		) VALUES (?, (SELECT provider_id from providers WHERE user_id = ?), ?, ?, ?)
	`,
		portfolio.ID,
		portfolio.UserID,
		portfolio.Title,
		portfolio.BookingID,
		portfolio.ServiceID,
	); err != nil {
		return err
	}

	// Save photos information if present
	if portfolio.Photos != nil {
		for _, photoUrl := range portfolio.Photos {
			photo := model.Photo{
				OwnerID:     portfolio.UserID,
				Url:         photoUrl,
				PortfolioID: portfolio.ID.String(),
			}
			err := createPhoto(ctx, tx, photo)
			if err != nil {
				log.Println("failed creating photo:", err)
			}
		}
	}

	return nil
}

func (s *PortfolioService) ListPortfoliosByProviderId(ctx context.Context, userId string) ([]*app.PortfolioBrief, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	portfolio, err := listPortfolioByCriteria(ctx, tx, "owner_id", userId)
	if err != nil {
		return nil, err
	}
	return portfolio, tx.Commit()
}

func (s *PortfolioService) ListPortfoliosByUserId(ctx context.Context, userId string) ([]*app.PortfolioBrief, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	provider, err := getProviderByCriteria(ctx, tx, "user_id", userId)
	if err != nil {
		return nil, err
	}

	portfolio, err := listPortfolioByCriteria(ctx, tx, "owner_id", provider.ID)
	if err != nil {
		return nil, err
	}
	/*
		portfolioWithPhotos, err := s.retrievePhotos(ctx, portfolio)
		if err == nil {
			portfolio = portfolioWithPhotos
		}
	*/
	return portfolio, tx.Commit()
}

//func listPortfolio(ctx context.Context, tx *Tx, userId string) ([]*app.Portfolio, error) {
func listPortfolioByCriteria(ctx context.Context, tx *Tx, haystack string, needle string) ([]*app.PortfolioBrief, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			portfolio_id,
			title
		FROM portfolios
		WHERE `+haystack+` = ?
	`, needle)
	if err != nil {
		log.Println("QueryCtx failed:", err)
		return nil, err
	}
	portfolios := make([]*app.PortfolioBrief, 0)
	defer rows.Close()
	for rows.Next() {
		portfolio := &app.PortfolioBrief{}
		err := rows.Scan(
			&portfolio.ID,
			&portfolio.Title,
		)
		if err != nil {
			log.Println("Scan failed:", err)
			return nil, err
		}
		/*
			res := retrievePhotos(tx, portfolio.ID.String())
			if res != nil {
				portfolio.Photos = res
			}
		*/
		portfolios = append(portfolios, portfolio)
	}

	return portfolios, nil
}

func (s *PortfolioService) retrievePhotos(ctx context.Context, portfolios []*app.Portfolio) ([]*app.Portfolio, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	portfolio := retrievePhotos(ctx, tx, portfolios)
	if portfolio == nil {
		return nil, errors.New("failed retrieving photos")
	}
	return portfolio, tx.Commit()
}

func retrievePhotos(ctx context.Context, tx *Tx, portfolios []*app.Portfolio) []*app.Portfolio {

	for _, portfolio := range portfolios {
		rows, err := tx.QueryContext(ctx, `
			SELECT
				photo_url
			FROM photos
			WHERE portfolio_id = ?
		`, portfolio.ID)

		if err != nil {
			log.Println("QueryCtx failed:", err)
			return nil
		}
		defer rows.Close()
		for rows.Next() {
			var photoUrl string
			if err := rows.Scan(
				&photoUrl,
			); err != nil {
				log.Println("Scan failed:", err)
				return nil
			}
			portfolio.Photos = append(portfolio.Photos, photoUrl)
		}
	}

	return portfolios
}

/*
	var photos []string
	// Get photos
	rows, err := tx.Query(`
	SELECT
		photo_url
	FROM photos
	WHERE portfolio_id = ?
	`, portfolioId)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var photo string
		if err := rows.Scan(
			&photo,
		); err != nil {
			return nil
		}
		photos = append(photos, photo)
	}
	return photos
}
*/

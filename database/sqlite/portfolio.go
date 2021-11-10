package sqlite

import (
	"context"

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
	err := tx.QueryRowContext(ctx, `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM portfolios
		WHERE id = ?
	`, id).Scan(
		&portfolio.ID,
		&portfolio.Title,
		&portfolio.BookingID,
	)
	if err != nil {
		return nil, err
	}
	return portfolio, nil
}

func (s *PortfolioService) CreatePortfolio(ctx context.Context, portfolio *model.Portfolio) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create dial and attach associated owner user.
	if err := createPortfolio(ctx, tx, portfolio); err != nil {
		return err
	}
	return tx.Commit()
}

func createPortfolio(ctx context.Context, tx *Tx, portfolio *model.Portfolio) error {
	// Create dial and attach associated owner user.
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO portfolio (
			id,
			name,
			booking_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
		)`,
		portfolio.ID,
		portfolio.Title,
		portfolio.BookingID,
	); err != nil {
		return err
	}
	return nil
}

func (s *PortfolioService) ListPortfoliosByUserId(ctx context.Context, userId string) ([]*app.Portfolio, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	portfolio, err := listPortfolio(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return portfolio, tx.Commit()
}

func listPortfolio(ctx context.Context, tx *Tx, userId string) ([]*app.Portfolio, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			portfolio.portfolio_id,
			portfolio.user_id,
			portfolio.service_id
		FROM portfolio
		WHERE portfolio.user_id = ?
		ORDER BY portfolio.created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	var portfolios []*app.Portfolio
	for rows.Next() {
		var portfolio app.Portfolio
		if err := rows.Scan(
			&portfolio.ID,
			&portfolio.Title,
			&portfolio.BookingID,
			&portfolio.Category,
		); err != nil {
			return nil, err
		}
		portfolios = append(portfolios, &portfolio)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return portfolios, nil
}

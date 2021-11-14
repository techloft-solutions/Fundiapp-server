package sqlite

import (
	"context"
	"log"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
)

type CategoryService struct {
	db *DB
}

func NewCategoryService(db *DB) *CategoryService {
	return &CategoryService{db}
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*app.Category, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	locations, err := retrieveCategories(ctx, tx)
	if err != nil {
		return nil, err
	}
	return locations, tx.Commit()
}

func retrieveCategories(ctx context.Context, tx *Tx) ([]*app.Category, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
			name,
			parent_id
		FROM categories
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over rows and deserialize into Dial objects.
	categories := make([]*app.Category, 0)
	for rows.Next() {
		var category app.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.ParentID,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *model.Category) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = createCategory(ctx, tx, category)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func createCategory(ctx context.Context, tx *Tx, category *model.Category) error {
	query := `
	INSERT INTO categories (
		name,
		profession,
		parent_id,
		description
	) VALUES (?, ?, ?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		category.Name,
		category.Profession,
		category.ParentID,
		category.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

type ReviewService struct {
	db *DB
}

func NewReviewService(db *DB) *ReviewService {
	return &ReviewService{db}
}

func (s *ReviewService) CreateReview(ctx context.Context, review *model.Review) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = createReview(ctx, tx, review)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func createReview(ctx context.Context, tx *Tx, review *model.Review) error {
	query := `
	INSERT INTO reviews (
		author_id,
		comment,
		rating,
		rating_quality,
		rating_resposiveness,
		rating_integrity,
		rating_competence,
		service_id
	) VALUES (?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		review.AuthorID,
		review.Comment,
		review.Rating,
		review.QualityRating,
		review.ResponsivenessRating,
		review.IntegrityRating,
		review.CompetenceRating,
		review.ServiceID,
	)
	if err != nil {
		return err
	}

	return nil
}

type SvcService struct {
	db *DB
}

func NewSvcService(db *DB) *ReviewService {
	return &ReviewService{db}
}

func (s *ReviewService) CreateService(ctx context.Context, service *model.Service) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = createService(ctx, tx, service)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func createService(ctx context.Context, tx *Tx, service *model.Service) error {
	query := `
	INSERT INTO services (
		user_id,
		provider_id,
		name,
		price,
		currency,
		price_unit
	) VALUES (?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		service.UserID,
		service.ProviderID,
		service.Name,
		service.Rate.Amount,
		service.Currency,
		service.Rate.Unit,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *LocationService) ListMyServices(ctx context.Context, userId string) ([]*app.Service, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	services, err := getServicesByUserID(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return services, tx.Commit()
}

func getServicesByUserID(ctx context.Context, tx *Tx, userId string) ([]*app.Service, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			name,
			price,
			currency,
			price_unit
		FROM services
		WHERE provider_id = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	services := make([]*app.Service, 0)
	for rows.Next() {
		var service app.Service
		if err := rows.Scan(
			&service.Name,
			&service.Rate.Price,
			&service.Rate.Currency,
			&service.Rate.Unit,
		); err != nil {
			log.Println("rows scan error:", err)
			return nil, err
		}
		services = append(services, &service)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

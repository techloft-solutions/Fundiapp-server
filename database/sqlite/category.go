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

func (s *CategoryService) ListRootCategories(ctx context.Context) ([]app.RootCategory, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	categories, err := retrieveRootCategories(ctx, tx)
	if err != nil {
		return nil, err
	}
	return categories, tx.Commit()
}

func retrieveRootCategories(ctx context.Context, tx *Tx) ([]app.RootCategory, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
			name
		FROM categories
		WHERE level = 0
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]app.RootCategory, 0)
	for rows.Next() {
		var category app.RootCategory
		if err := rows.Scan(
			&category.ID,
			&category.Name,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*app.Category, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	categories, err := retrieveCategoriesByCriteria(ctx, tx, "1", "1")
	if err != nil {
		return nil, err
	}
	return categories, tx.Commit()
}

func (s *CategoryService) ListCategoriesByParentID(ctx context.Context, parentID string) ([]*app.Category, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	categories, err := retrieveCategoriesByCriteria(ctx, tx, "parent_id", parentID)
	if err != nil {
		return nil, err
	}
	return categories, tx.Commit()
}

func (s *CategoryService) ListCategoriesByIndustryID(ctx context.Context, industryID string) ([]*app.Category, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	categories, err := retrieveCategoriesByCriteria(ctx, tx, "industry_id", industryID)
	if err != nil {
		return nil, err
	}
	return categories, tx.Commit()
}

func retrieveCategoriesByCriteria(ctx context.Context, tx *Tx, haystack string, needle string) ([]*app.Category, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
			name,
			parent_id,
			icon_url
		FROM categories
		WHERE `+haystack+` = ?
		AND level > 0
		`,
		needle,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*app.Category, 0)
	for rows.Next() {
		var category app.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.ParentID,
			&category.IconURL,
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
		parent_id,
		description,
		icon_url
	) VALUES (?, ?, ?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		category.Name,
		category.ParentID,
		category.Description,
		category.IconURL,
	)
	if err != nil {
		return err
	}

	return nil
}

type IndustryService struct {
	db *DB
}

func NewIndustryService(db *DB) *IndustryService {
	return &IndustryService{db}
}

func (s *IndustryService) CreateIndustry(ctx context.Context, industry *model.Industry) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = createIndustry(ctx, tx, industry)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func createIndustry(ctx context.Context, tx *Tx, industry *model.Industry) error {
	query := `
	INSERT INTO industries (
		name,
		description,
		icon_url
	) VALUES (?, ?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		industry.Name,
		industry.Description,
		industry.IconURL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *IndustryService) ListIndustries(ctx context.Context) ([]*app.Industry, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	industries, err := retrieveIndustries(ctx, tx)
	if err != nil {
		return nil, err
	}
	return industries, tx.Commit()
}

func retrieveIndustries(ctx context.Context, tx *Tx) ([]*app.Industry, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
			name,
			icon_url
		FROM industries
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	industries := make([]*app.Industry, 0)
	for rows.Next() {
		var industry app.Industry
		if err := rows.Scan(
			&industry.ID,
			&industry.Name,
			&industry.IconURL,
		); err != nil {
			return nil, err
		}
		industries = append(industries, &industry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return industries, nil
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
		provider_id,
		comment,
		rating,
		rating_quality,
		rating_resposiveness,
		rating_integrity,
		rating_competence,
		service_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		review.AuthorID,
		review.ProviderID,
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

func (s *ReviewService) ListReviewsByProviderID(ctx context.Context, providerId string) ([]*app.Review, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	reviews, err := getReviewsByProviderID(ctx, tx, providerId)
	if err != nil {
		return nil, err
	}
	return reviews, tx.Commit()
}

func getReviewsByProviderID(ctx context.Context, tx *Tx, providerId string) ([]*app.Review, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			author_id,
			comment,
			rating,
			rating_quality,
			rating_resposiveness,
			rating_integrity,
			rating_competence,
			service_id
		FROM reviews
		WHERE provider_id = ?
		`,
		providerId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*app.Review, 0)
	for rows.Next() {
		var review app.Review
		if err := rows.Scan(
			&review.ID,
			//&review.AuthorID,
			&review.Comment,
			&review.Rating,
			//&review.QualityRating,
			//&review.ResponsivenessRating,
			//&review.IntegrityRating,
			//&review.CompetenceRating,
			//&review.ServiceID,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *UserService) CreateService(ctx context.Context, service *model.Service) error {
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
		provider_id,
		name,
		price,
		currency,
		price_unit
	) VALUES (?, ?, ?, ?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
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

func (s *UserService) ListMyServices(ctx context.Context, userId string) ([]*app.Service, error) {
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

func (s *UserService) ListServicesByProviderID(ctx context.Context, providerId string) ([]*app.Service, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	services, err := getServicesByProviderID(ctx, tx, providerId)
	if err != nil {
		return nil, err
	}
	return services, tx.Commit()
}

func getServicesByProviderID(ctx context.Context, tx *Tx, providerId string) ([]*app.Service, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			name,
			price,
			currency,
			price_unit
		FROM services
		WHERE provider_id = ?
	`, providerId)
	if err != nil {
		return nil, err
	}
	services := make([]*app.Service, 0)
	for rows.Next() {
		var service app.Service
		if err := rows.Scan(
			&service.ID,
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

func getServicesByUserID(ctx context.Context, tx *Tx, userId string) ([]*app.Service, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			name,
			price,
			currency,
			price_unit
		FROM services
		WHERE provider_id IN (
			SELECT provider_id FROM users WHERE user_id = ?
		)
	`, userId)
	if err != nil {
		return nil, err
	}

	services := make([]*app.Service, 0)
	for rows.Next() {
		var service app.Service
		if err := rows.Scan(
			&service.ID,
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

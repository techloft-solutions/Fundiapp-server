package sqlite

import (
	"context"

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
	INSERT INTO category (
		title,
		profession,
		parent_id,
		description
	) VALUES (?, ?)
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

func retrieveCategories(ctx context.Context, tx *Tx) ([]*app.Category, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
			title,
			profession
		FROM categories
		ORDER BY start_date ASC
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

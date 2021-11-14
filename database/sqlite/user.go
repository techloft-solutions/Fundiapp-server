package sqlite

import (
	"context"
	"database/sql"
	"log"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

//func (s *UserService) CreateProvider(ctx context.Context, booking *model.Provider) error {}

//func (s *UserService) CreateClient(ctx context.Context, booking *model.Client) error {}

//func (s *UserService) FindClientByID(ctx context.Context, booking *model.Client) error {}

func (s *UserService) FindProviderByID(ctx context.Context, id string) (*app.Provider, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	profile, err := getProviderProfileByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return profile, tx.Commit()
}

func getProviderProfileByID(ctx context.Context, tx *Tx, id string) (*app.Provider, error) {
	profile := &app.Provider{}
	err := tx.QueryRowContext(ctx, `
		SELECT
			providers.provider_id,
			providers.user_id,
			profiles.first_name,
			profiles.last_name,
			providers.bio,
			providers.profession,
			providers.ratings_average,
			providers.reviews_count,
			providers.services_count,
			providers.portfolio_count,
			locations.name
		FROM providers
		LEFT JOIN profiles ON profiles.user_id = providers.user_id
		LEFT JOIN locations ON profiles.location_id = locations.location_id
		WHERE providers.provider_id = ?
	`, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Bio,
		&profile.Profession,
		&profile.AvgRating,
		&profile.Stats.Reviews,
		&profile.Stats.Services,
		&profile.Stats.Portfolios,
		&profile.Location.Name,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *UserService) CreateProvider(ctx context.Context, provider *model.Provider) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createProvider(ctx, tx, provider); err != nil {
		return err
	}

	log.Println("Provider created")

	if err := createProfile(ctx, tx, &provider.Profile); err != nil {
		return err
	}

	return tx.Commit()
}

func createProvider(ctx context.Context, tx *Tx, provider *model.Provider) error {
	query := `
	INSERT INTO providers (
		provider_id,
		user_id,
		bio,
		profession
	) VALUES (?, ?, ?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		provider.ID,
		provider.UserID,
		provider.Bio,
		provider.Profession,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ListProviders(ctx context.Context) ([]*app.ProviderBrief, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	providers, err := findProviders(ctx, tx)
	if err != nil {
		return nil, err
	}

	return providers, nil
}

func findProviders(ctx context.Context, tx *Tx) ([]*app.ProviderBrief, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    providers.provider_id,
			profiles.first_name,
			profiles.last_name,
			providers.profession,
			providers.ratings_average,
			providers.reviews_count,
			providers.jobs_count,
			providers.rate_per_hour,
			providers.currency
		FROM providers
		LEFT JOIN profiles ON profiles.user_id = providers.user_id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over rows and deserialize into Dial objects.
	providers := make([]*app.ProviderBrief, 0)
	for rows.Next() {
		var fname string
		var lname string
		var provider app.ProviderBrief
		if err := rows.Scan(
			&provider.ID,
			&fname,
			&lname,
			&provider.Profession,
			&provider.Rating,
			&provider.Reviews,
			&provider.Jobs,
			&provider.Rate.Price,
			&provider.Currency,
		); err != nil {
			return nil, err
		}
		provider.Name = fname + " " + lname
		providers = append(providers, &provider)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return providers, nil
}

func (s *UserService) GetProfile(ctx context.Context, userId string) (*app.Profile, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	profile, err := getProfileByUserID(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return profile, tx.Commit()
}

func (s *UserService) FindProfileByUserID(ctx context.Context, userId string) (*app.Profile, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	profile, err := getProfileByUserID(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return profile, tx.Commit()
}

func getProfileByUserID(ctx context.Context, tx *Tx, userId string) (*app.Profile, error) {
	profile := &app.Profile{}
	profile.UserID = userId
	err := tx.QueryRowContext(ctx, `
		SELECT
			profile_id,
			first_name,
			last_name,
			location_id,
			verified
		FROM profiles
		WHERE user_id = ?
	`, userId).Scan(
		&profile.ID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Location.ID,
		&profile.Verified,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, profile *model.Profile) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updateProfile(ctx, tx, profile); err != nil {
		return err
	}
	return tx.Commit()
}

func updateProfile(ctx context.Context, tx *Tx, profile *model.Profile) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE profiles as p
		SET first_name = COALESCE(?, first_name),
			last_name = COALESCE(?, last_name),
			location_id = COALESCE(?, location_id)
		WHERE user_id = ?
	`,
		profile.FirstName,
		profile.LastName,
		profile.LocationID,
		profile.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *UserService) CreateProfile(ctx context.Context, profile *model.Profile) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createProfile(ctx, tx, profile); err != nil {
		return err
	}
	return tx.Commit()
}

func createProfile(ctx context.Context, tx *Tx, profile *model.Profile) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO profiles (
			profile_id,
			user_id, 
			first_name, 
			last_name, 
			location_id,
			account_type
		) VALUES (?, ?, ?, ?, ?, ?)
	`,
		profile.ID,
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.LocationID,
		profile.Type,
	)
	if err != nil {
		return err
	}
	return nil
}

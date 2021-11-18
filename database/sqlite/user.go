package sqlite

import (
	"context"
	"database/sql"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/google/uuid"
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

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}
	return tx.Commit()
}

func createUser(ctx context.Context, tx *Tx, user *model.User) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO users (
			user_id,
			username,
			phone,
			password,
			is_provider
		) VALUES (?, ?, ?, ?, ?)
		`,
		user.UserID,
		user.Username,
		user.Phone,
		user.Password,
		user.IsProvider,
	)
	if err != nil {
		return err
	}

	if user.IsProvider {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO providers (
				provider_id,
				user_id
			) VALUES (?, ?)
			`,
			uuid.New(),
			user.UserID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, user *model.ResetUser) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updateUserPassword(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func updateUserPassword(ctx context.Context, tx *Tx, user *model.ResetUser) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE users SET
			password = ?
		WHERE phone = ? AND reset_password_code = ?
		`,
		user.NewPassword,
		user.Phone,
		user.ResetCode,
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

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updateUser(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func updateUser(ctx context.Context, tx *Tx, user *model.User) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE users SET
			username = COALESCE(?, username),
			reset_password_code = COALESCE(?, reset_password_code)
		WHERE id = ?
		`,
		user.Username,
		user.ResetCode,
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

func (s *UserService) FindProviderByUserID(ctx context.Context, userId string) (*app.Provider, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	profile, err := getProviderByUserID(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return profile, tx.Commit()
}

func (s *UserService) UpdateResetCode(ctx context.Context, code int, phone string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updateResetCode(ctx, tx, code, phone); err != nil {
		return err
	}

	return tx.Commit()
}

func updateResetCode(ctx context.Context, tx *Tx, code int, phone string) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE users SET
			reset_password_code = ?
		WHERE phone = ?
		`,
		code,
		phone,
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

func getProviderByUserID(ctx context.Context, tx *Tx, id string) (*app.Provider, error) {
	provider := &app.Provider{}
	err := tx.QueryRowContext(ctx, `
		SELECT
			providers.provider_id
		FROM providers
		WHERE providers.user_id = ?
	`, id).Scan(
		&provider.ID,
	)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

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
			users.first_name,
			users.last_name,
			providers.bio,
			providers.profession,
			providers.ratings_average,
			providers.reviews_count,
			providers.services_count,
			providers.portfolio_count,
			locations.name
		FROM providers
		LEFT JOIN users ON users.user_id = providers.user_id
		LEFT JOIN locations ON locations.location_id = users.location_id
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
		&profile.Location,
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
	/*
		if err := createProfile(ctx, tx, &provider.Profile); err != nil {
			return err
		}
	*/
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
			users.user_id,
			users.username,
			providers.profession,
			providers.ratings_average,
			providers.reviews_count,
			providers.jobs_count,
			providers.rate_per_hour,
			providers.currency,
			users.photo_url
		FROM providers
		LEFT JOIN users ON users.user_id = providers.user_id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := make([]*app.ProviderBrief, 0)
	for rows.Next() {
		var provider app.ProviderBrief
		if err := rows.Scan(
			&provider.ID,
			&provider.UserID,
			&provider.Name,
			&provider.Profession,
			&provider.Rating,
			&provider.Reviews,
			&provider.Jobs,
			&provider.Rate.Price,
			&provider.Currency,
			&provider.Photo,
		); err != nil {
			return nil, err
		}
		providers = append(providers, &provider)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return providers, nil
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
			p.username,
			p.first_name,
			p.last_name,
			p.email,
			p.photo_url,
			locations.address,
			p.verified
		FROM users as p
		LEFT JOIN locations ON locations.location_id = p.location_id
		WHERE p.user_id = ? AND p.is_provider = 0
	`, userId).Scan(
		&profile.Username,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.PhotoUrl,
		&profile.Location,
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
		UPDATE users as p
		SET
			first_name = COALESCE(?, first_name),
			last_name = COALESCE(?, last_name),
			email = COALESCE(?, email),
			photo_url = COALESCE(?, photo_url),
			location_id = COALESCE(?, location_id)
		WHERE user_id = ?
	`,
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.PhotoUrl,
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

func (s *UserService) UpdateProvider(ctx context.Context, provider *model.Provider) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updateProvider(ctx, tx, provider); err != nil {
		return err
	}

	if err := updateProfile(ctx, tx, &provider.Profile); err != nil {
		return err
	}
	return tx.Commit()
}

func updateProvider(ctx context.Context, tx *Tx, provider *model.Provider) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE providers
		SET
			bio = COALESCE(?, bio),
			profession = COALESCE(?, profession)
		WHERE user_id = ? 
	`,
		provider.Bio,
		provider.Profession,
		provider.UserID,
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
		INSERT INTO users (
			user_id, 
			first_name, 
			last_name,
			email,
			photo_url,
			location_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.PhotoUrl,
		profile.LocationID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) FindUserByUsername(ctx context.Context, username string) (*app.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := getUserByCriteria(ctx, tx, "username", username)
	if err != nil {
		return nil, err
	}
	return user, tx.Commit()
}

func (s *UserService) FindUserByPhoneNumber(ctx context.Context, phone string) (*app.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := getUserByCriteria(ctx, tx, "phone", phone)
	if err != nil {
		return nil, err
	}
	return user, tx.Commit()
}

func (s *UserService) FindUserByID(ctx context.Context, userID string) (*app.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := getUserByCriteria(ctx, tx, "user_id", userID)
	if err != nil {
		return nil, err
	}
	return user, tx.Commit()
}

func getUserByCriteria(ctx context.Context, tx *Tx, haystack string, needle string) (*app.User, error) {
	user := &app.User{}
	err := tx.QueryRowContext(ctx, `
		SELECT
			user_id,
			username,
			first_name,
			last_name,
			email,
			photo_url,
			phone,
			is_provider			
		FROM users
		WHERE `+haystack+` = ?
	`, needle).Scan(
		&user.UserID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhotoUrl,
		&user.Phone,
		&user.IsProvider,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ValidateUser(ctx context.Context, phone string, password string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := validateUser(ctx, tx, phone, password); err != nil {
		return err
	}
	return tx.Commit()
}

func validateUser(ctx context.Context, tx *Tx, phone string, password string) error {
	err := tx.QueryRowContext(ctx, `
		SELECT
			username
		FROM users
		WHERE phone = ? AND password = ?
	`, phone, password).Scan(
		&phone,
	)
	if err != nil {
		return err
	}
	return nil
}

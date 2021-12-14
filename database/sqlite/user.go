package sqlite

import (
	"context"
	"database/sql"
	"log"
	"strings"

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

func (s *UserService) ResetUserPassword(ctx context.Context, password string, userID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := resetUserPassword(ctx, tx, password, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func resetUserPassword(ctx context.Context, tx *Tx, newPass string, userID string) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE users SET
			password = ?
		WHERE user_id = ?
		`,
		newPass,
		userID,
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

func (s *UserService) ChangeUserPassword(ctx context.Context, user *model.PwdChange) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := changePassword(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func changePassword(ctx context.Context, tx *Tx, user *model.PwdChange) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE users SET
			password = ?
		WHERE user_id = ? AND password = ?
		`,
		user.NewPassword,
		user.UserID,
		user.OldPassword,
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

	provider, err := getProviderByCriteria(ctx, tx, "user_id", userId)
	if err != nil {
		return nil, err
	}

	service, err := getServicesByProviderID(ctx, tx, provider.ID)
	if err != nil {
		return nil, err
	}
	provider.Services = service
	return provider, tx.Commit()
}

func getProviderByCriteria(ctx context.Context, tx *Tx, haystack string, needle string) (*app.Provider, error) {
	provider := &app.Provider{}
	location := app.ProfileLocation{}
	price := app.Price{}
	err := tx.QueryRowContext(ctx, `
		SELECT
			providers.provider_id,
			providers.user_id,
			users.username,
			users.first_name,
			users.last_name,
			users.phone,
			users.email,
			users.photo_url,
			providers.bio,
			categories.name AS profession,
			providers.ratings_average,
			providers.reviews_count,
			providers.services_count,
			providers.portfolio_count,
			providers.price,
			providers.currency,
			locations.location_id,
			locations.name
		FROM providers
		LEFT JOIN users ON users.user_id = providers.user_id
		LEFT JOIN locations ON locations.location_id = users.location_id
		LEFT JOIN categories ON categories.id = providers.category_id
		WHERE providers.`+haystack+` = ?
	`, needle).Scan(
		&provider.ID,
		&provider.UserID,
		&provider.Username,
		&provider.FirstName,
		&provider.LastName,
		&provider.Phone,
		&provider.Email,
		&provider.PhotoUrl,
		&provider.Bio,
		&provider.Profession,
		&provider.AvgRating,
		&provider.Stats.Reviews,
		&provider.Stats.Services,
		&provider.Stats.Portfolios,
		&price.Amount,
		&price.Currency,
		&location.ID,
		&location.Address,
	)
	if err != nil {
		return nil, err
	}

	if location != (app.ProfileLocation{}) {
		provider.Location = &location
	}

	if price.Amount != nil {
		provider.Price = &price
	}

	return provider, nil
}

func (s *UserService) FindProviderByID(ctx context.Context, id string) (*app.Provider, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	provider, err := getProviderByCriteria(ctx, tx, "provider_id", id)
	if err != nil {
		return nil, err
	}

	service, err := getServicesByProviderID(ctx, tx, provider.ID)
	if err != nil {
		return nil, err
	}
	provider.Services = service

	return provider, tx.Commit()
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

func (s *UserService) FilterProviders(ctx context.Context, filter model.ProviderFilter) ([]*app.ProviderBrief, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	providers, err := filterProviders(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	return providers, nil
}

func findProviders(ctx context.Context, tx *Tx) ([]*app.ProviderBrief, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    providers.provider_id,
			CONCAT_WS(' ', users.first_name, users.last_name) AS full_name,
			categories.name AS profession,
			providers.ratings_average,
			providers.reviews_count,
			providers.jobs_count,
			providers.price,
			providers.currency,
			users.photo_url
		FROM providers
		INNER JOIN users ON users.user_id = providers.user_id
		LEFT JOIN categories ON categories.id = providers.category_id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := make([]*app.ProviderBrief, 0)
	for rows.Next() {
		var provider app.ProviderBrief
		price := app.Price{}
		if err := rows.Scan(
			&provider.ID,
			&provider.Name,
			&provider.Profession,
			&provider.Rating,
			&provider.Reviews,
			&provider.Jobs,
			&price.Amount,
			&price.Currency,
			&provider.Photo,
		); err != nil {
			return nil, err
		}
		if price != (app.Price{}) {
			provider.Price = &price
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
	location := app.ProfileLocation{}
	err := tx.QueryRowContext(ctx, `
		SELECT
			p.username,
			p.first_name,
			p.last_name,
			p.email,
			p.photo_url,
			locations.location_id,
			locations.address,
			p.verified
		FROM users as p
		LEFT JOIN locations ON locations.location_id = p.location_id
		WHERE p.user_id = ?
	`, userId).Scan(
		&profile.Username,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.PhotoUrl,
		&location.ID,
		&location.Address,
		&profile.Verified,
	)
	if err != nil {
		return nil, err
	}
	if location != (app.ProfileLocation{}) {
		profile.Location = &location
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
		UPDATE users
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

	success := true

	if errProvider := updateProvider(ctx, tx, provider); err != nil {
		log.Println("error updating provider:", errProvider)
		success = false
	}

	if errProfile := updateProfile(ctx, tx, &provider.Profile); err != nil {
		log.Println("error updating profile:", errProfile)
		if !success {
			return errProfile
		}
	}

	return tx.Commit()
}

func updateProvider(ctx context.Context, tx *Tx, provider *model.Provider) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE providers
		SET
			bio = COALESCE(?, bio),
			category_id = COALESCE(?, category_id),
			industry_id = COALESCE(?, industry_id)
		WHERE user_id = ? 
	`,
		provider.Bio,
		provider.CategoryID,
		provider.IndustryID,
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

func (s *UserService) ValidateUser(ctx context.Context, phone string, password string, isProvider bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if isProvider {
		if err := validateUserAsProvider(ctx, tx, phone, password); err != nil {
			return err
		}
	} else {
		if err := validateUser(ctx, tx, phone, password); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func validateUser(ctx context.Context, tx *Tx, phone string, password string) error {
	err := tx.QueryRowContext(ctx, `
		SELECT
			username
		FROM users
		WHERE phone = ? AND password = ? AND is_provider = false
	`, phone, password).Scan(
		&phone,
	)
	if err != nil {
		return err
	}
	return nil
}

func validateUserAsProvider(ctx context.Context, tx *Tx, phone string, password string) error {
	err := tx.QueryRowContext(ctx, `
		SELECT
			username
		FROM users
		WHERE phone = ? AND password = ? AND is_provider = true
	`, phone, password).Scan(
		&phone,
	)
	if err != nil {
		return err
	}
	return nil
}

func filterProviders(ctx context.Context, tx *Tx, filter model.ProviderFilter) (_ []*app.ProviderBrief, err error) {
	// Build WHERE clause. Each part of the WHERE clause is AND-ed together.
	// Values are appended to an arg list to avoid SQL injection.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.IndustryID; v != "" {
		where, args = append(where, "industry_id = ?"), append(args, v)
	}

	if v := filter.CategoryID; v != "" {
		where, args = append(where, "category_id = ?"), append(args, v)
	}

	// Execue query with limiting WHERE clause and LIMIT/OFFSET injected.
	rows, err := tx.QueryContext(ctx, `
	SELECT 
			providers.provider_id,
			CONCAT_WS(' ', users.first_name, users.last_name) AS full_name,
			categories.name AS profession,
			providers.ratings_average,
			providers.reviews_count,
			providers.jobs_count,
			providers.price,
			providers.currency,
			users.photo_url
		FROM providers
		INNER JOIN users ON users.user_id = providers.user_id
		LEFT JOIN categories ON categories.id = providers.category_id
		WHERE `+strings.Join(where, " AND ")+`
		AND providers.is_verfied = true
		ORDER BY providers.updated_at ASC
		`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := make([]*app.ProviderBrief, 0)
	for rows.Next() {
		var provider app.ProviderBrief
		//var price app.Price
		if err := rows.Scan(
			&provider.ID,
			&provider.Name,
			&provider.Profession,
			&provider.Rating,
			&provider.Reviews,
			&provider.Jobs,
			&provider.Price.Amount,
			&provider.Price.Currency,
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

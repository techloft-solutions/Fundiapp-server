package sqlite

import (
	"context"
	"database/sql"
	"log"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/google/uuid"
)

type BookingService struct {
	db *DB
}

func NewBookingService(db *DB) *BookingService {
	return &BookingService{db}
}

type RequestService struct {
	db *DB
}

func NewRequestService(db *DB) *RequestService {
	return &RequestService{db}
}

func (s *RequestService) CreateRequest(ctx context.Context, request *model.Request) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createRequest(ctx, tx, request); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RequestService) ListRequests(ctx context.Context, userId app.UserID) ([]*app.Request, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	requests, err := listRequestsByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return requests, tx.Commit()
}

func (s *RequestService) FindRequestByID(ctx context.Context, id uuid.UUID) (*app.RequestDetail, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	request, err := findRequestByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return request, tx.Commit()
}

func findRequestByID(ctx context.Context, tx *Tx, id uuid.UUID) (*app.RequestDetail, error) {
	request := &app.RequestDetail{}
	err := tx.QueryRowContext(ctx, `
		SELECT
			booking_id,
			title,
			description,
			client_id,
			start_date,
			created_at
		FROM bookings
		WHERE booking_id = ?
	`, id).Scan(
		&request.ID,
		&request.Title,
		&request.Note,
		&request.Client,
		&request.Start,
		&request.Created,
	)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func listRequestsByUserId(ctx context.Context, tx *Tx, userId app.UserID) ([]*app.Request, error) {
	requests := []*app.Request{}
	rows, err := tx.QueryContext(ctx, `
		SELECT
			booking_id,
			title,
			start_date,
			created_at
		FROM bookings
		WHERE client_id = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		request := &app.Request{}
		if err := rows.Scan(
			&request.ID,
			&request.Title,
			&request.Start,
			&request.Created,
		); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func createRequest(ctx context.Context, tx *Tx, request *model.Request) error {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO bookings (
			booking_id,
			client_id,
			title,
			description,
			start_date,
			location_id,
			status
		) VALUES (?,?,?,?,?,?,?)
		`,
		request.ID,
		request.ClientID,
		request.Title,
		request.Note,
		request.StartDate,
		request.LocationID,
		request.Status,
	); err != nil {
		panic(err)
	}

	// Save photos information if present
	if request.Photos != nil {
		for _, photoUrl := range request.Photos {
			photo := model.Photo{
				Owner: request.ClientID,
				Url:   photoUrl,
			}
			photo.ID = uuid.New()
			err := createPhoto(ctx, tx, photo)
			if err != nil {
				log.Println(err)
			}
			bookingPhoto := model.BookingPhoto{
				BookingID: request.ID,
				PhotoID:   photo.ID,
			}
			createBookingPhoto(ctx, tx, bookingPhoto)
		}
	}
	return nil
}

type LocationService struct {
	db *DB
}

func NewLocationService(db *DB) *LocationService {
	return &LocationService{db}
}

func (s *LocationService) ListMyLocations(ctx context.Context, authUser *app.AuthUser) ([]*app.Location, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	locations, err := getLocationsByUserID(ctx, tx, authUser.ID)
	if err != nil {
		return nil, err
	}
	return locations, tx.Commit()
}

func (s *LocationService) CreateLocation(ctx context.Context, location *model.Location) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create location and attach associated owner user.
	if err := createLocation(ctx, tx, location); err != nil {
		return err
	}
	return tx.Commit()
}

func createLocation(ctx context.Context, tx *Tx, location *model.Location) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO locations(
		location_id,
		user_id,
		title,
		address,
		latitude,
		longitude,
		) VALUES (?,?,?,?,?,?,?,?)
		`,
		location.ID,
		location.UserID,
		location.Title,
		location.Address,
		location.Latitude,
		location.Longitude,
	)
	if err != nil {
		return err
	}
	return nil
}

func getLocationsByUserID(ctx context.Context, tx *Tx, userId string) ([]*app.Location, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			locations.location_id,
			locations.name,
			locations.address,
			locations.latitude,
			locations.longitude
		FROM locations
		WHERE locations.user_id = ?
		ORDER BY locations.created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	var locations []*app.Location
	for rows.Next() {
		var location app.Location
		if err := rows.Scan(
			&location.ID,
			&location.Title,
			&location.Address,
			&location.Latitude,
			&location.Longitude,
		); err != nil {
			return nil, err
		}
		locations = append(locations, &location)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return locations, nil
}

type TransactionService struct {
	db *DB
}

func NewTransactionService(db *DB) *TransactionService {
	return &TransactionService{db}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create dial and attach associated owner user.
	if err := createTransaction(ctx, tx, transaction); err != nil {
		return err
	}
	return tx.Commit()
}

func createTransaction(ctx context.Context, tx *Tx, transaction *model.Transaction) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO transactions(
		transaction_id,
		user_id,
		code,
		service_id,
		amount,
		) VALUES (?,?,?,?,?,?,?)
		`,
		transaction.ID,
		transaction.UserID,
		transaction.Code,
		transaction.ServiceID,
		transaction.Amount,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionService) ListTransactions(ctx context.Context, userId uuid.UUID) ([]*app.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	transactions, err := listTransactions(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return transactions, tx.Commit()
}

func listTransactions(ctx context.Context, tx *Tx, userId uuid.UUID) ([]*app.Transaction, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			transactions.transaction_id,
			transactions.user_id,
			transactions.location_id,
			transactions.booking_id,
			transactions.amount,
			transactions.created_at,
			transactions.updated_at
		FROM transactions
		WHERE transactions.user_id = ?
		`, userId)
	if err != nil {
		return nil, err
	}
	var transactions []*app.Transaction
	for rows.Next() {
		var transaction app.Transaction
		if err := rows.Scan(
			&transaction.Code,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Amount,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *LocationService) FindLocationsByUserID(ctx context.Context, userId uuid.UUID) ([]*app.Location, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	locations, err := findLocationsByUserID(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return locations, tx.Commit()
}

func findLocationsByUserID(ctx context.Context, tx *Tx, userId uuid.UUID) ([]*app.Location, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			locations.location_id,
			locations.name,
			locations.address,
			locations.latitude,
			locations.longitude,
			locations.created_at,
			locations.updated_at
		FROM locations
		INNER JOIN users_locations
		ON locations.location_id = users_locations.location_id
		WHERE users_locations.user_id = ?
		`, userId)
	if err != nil {
		return nil, err
	}

	var locations []*app.Location
	for rows.Next() {
		var location app.Location
		if err := rows.Scan(
			&location.ID,
			&location.Title,
			&location.Address,
			&location.Latitude,
			&location.Longitude,
		); err != nil {
			return nil, err
		}
		locations = append(locations, &location)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return locations, nil
}

func (s *BookingService) FindBookingByID(ctx context.Context, id uuid.UUID) (_ *app.Booking, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create dial and attach associated owner user.
	booking, err := findBookingByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return booking, tx.Commit()
}

func (s *BookingService) FindBookings(ctx context.Context) ([]*app.BookingBrief, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	bookings, err := findBookings(ctx, tx)
	if err != nil {
		return nil, err
	}
	return bookings, tx.Commit()
}

func (s *BookingService) CreateBooking(ctx context.Context, booking *model.Booking) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createBooking(ctx, tx, booking); err != nil {
		return err
	}
	return tx.Commit()
}

func createPhoto(ctx context.Context, tx *Tx, photo model.Photo) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO photos(
		photo_id,
		uploaded_by,
		photo_url
		) VALUES (?,?,?)
		`,
		photo.ID,
		photo.Owner,
		photo.Url,
	)
	if err != nil {
		return err
	}
	return nil
}

func createBookingPhoto(ctx context.Context, tx *Tx, bp model.BookingPhoto) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO booking_photos(
		photo_id,
		booking_id
		) VALUES (?,?)
		`,
		bp.PhotoID,
		bp.BookingID,
	)
	if err != nil {
		return err
	}
	return nil
}

// createBooking creates a new booking.
func createBooking(ctx context.Context, tx *Tx, booking *model.Booking) error {
	booking.Status = statusPending

	query := `
	INSERT INTO bookings (
		booking_id,
		title,
		description,
		status,
		start_date,
		client_id,
		provider_id,
		location_id,
		service_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		booking.ID,
		booking.Title,
		booking.Description,
		booking.Status,
		booking.StartDate,
		booking.ClientID,
		booking.ProviderID,
		booking.LocationID,
		booking.ServiceID,
	)
	if err != nil {
		panic(err)
	}

	// Save photos information if present
	if booking.Photos != nil {
		for _, photoUrl := range booking.Photos {
			photo := model.Photo{
				Owner: booking.ClientID,
				Url:   photoUrl,
			}
			photo.ID = uuid.New()
			err := createPhoto(ctx, tx, photo)
			if err != nil {
				log.Println(err)
			}
			bookingPhoto := model.BookingPhoto{
				BookingID: booking.ID,
				PhotoID:   photo.ID,
			}
			createBookingPhoto(ctx, tx, bookingPhoto)
		}
	}

	return nil
}

func findBookings(ctx context.Context, tx *Tx) ([]*app.BookingBrief, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    booking_id,
			title,
			status,
			description,
			created_at,
			start_date
		FROM bookings
		ORDER BY start_date ASC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over rows and deserialize into Dial objects.
	bookings := make([]*app.BookingBrief, 0)
	for rows.Next() {
		var booking app.BookingBrief
		if err := rows.Scan(
			&booking.ID,
			&booking.Title,
			&booking.Status,
			&booking.Description,
			&booking.BookedAt,
			&booking.StartAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// findDials retrieves a list of matching dials. Also returns a total matching
// count which may different from the number of results if filter.Limit is set.
func findBookingByID(ctx context.Context, tx *Tx, id uuid.UUID) (_ *app.Booking, err error) {
	booking := &app.Booking{}
	// Execue query with limiting WHERE clause and LIMIT/OFFSET injected.
	if err := tx.QueryRowContext(ctx, `
		SELECT
			b.booking_id,
			b.title,
			b.description,
			b.status,
			b.start_date,
			b.client_id,
			b.provider_id,
			b.created_at
		FROM
			bookings b
		LEFT JOIN services s ON s.service_id = b.service_id
		WHERE
			b.booking_id = ?
		ORDER BY
			b.start_date ASC
		`,
		id,
	).Scan(&booking.ID, &booking.Title, &booking.Description, &booking.Status, &booking.StartAt, &booking.Client.UserID, &booking.Provider.UserID, &booking.BookedAt); err != nil {
		if err == sql.ErrNoRows {
			return booking, nil
		}
		return nil, err
	}

	return booking, nil
}

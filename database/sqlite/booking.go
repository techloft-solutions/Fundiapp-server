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

func createRequest(ctx context.Context, tx *Tx, request *model.Request) error {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO bookings (
			booking_id,
			client_id,
			title,
			description,
			start_at,
			location_id,
			status,
			is_urgent,
			is_request
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		request.ID,
		request.ClientID,
		request.Title,
		request.Note,
		request.StartDate,
		request.LocationID,
		request.Status,
		request.Urgent,
		true,
	); err != nil {
		log.Println("failed inserting bookings into db:", err)
		return err
	}

	// Save photos information if present
	if request.Photos != nil {
		for _, photoUrl := range request.Photos {
			photo := model.Photo{
				OwnerID:   request.ClientID,
				Url:       photoUrl,
				BookingID: request.ID.String(),
			}
			err := createPhoto(ctx, tx, photo)
			if err != nil {
				log.Println("failed creating photo:", err)
			}
		}
	}
	return nil
}

func (s *RequestService) ListRequests(ctx context.Context, userId app.UserID) ([]app.Request, error) {
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

func listRequestsByUserId(ctx context.Context, tx *Tx, userId app.UserID) ([]app.Request, error) {
	requests := []app.Request{}
	rows, err := tx.QueryContext(ctx, `
		SELECT
			booking_id,
			title,
			status,
			start_at,
			created_at
		FROM bookings
		WHERE client_id = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		request := app.Request{}
		if err := rows.Scan(
			&request.ID,
			&request.Title,
			&request.Status,
			&request.Start,
			&request.Created,
		); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
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
			status,
			start_at,
			created_at
		FROM bookings
		WHERE booking_id = ?
	`, id).Scan(
		&request.ID,
		&request.Title,
		&request.Note,
		&request.Client,
		&request.Status,
		&request.Start,
		&request.Created,
	)
	if err != nil {
		return nil, err
	}

	// Get photos
	rows, err := tx.QueryContext(ctx, `
		SELECT
			photo_url
		FROM photos
		WHERE booking_id = ?
	`, id)
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
		request.Photos = append(request.Photos, photo)
	}

	return request, nil
}

type LocationService struct {
	db *DB
}

func NewLocationService(db *DB) *LocationService {
	return &LocationService{db}
}

func (s *LocationService) ListMyLocations(ctx context.Context, userId string) ([]*app.Location, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	locations, err := getLocationsByUserID(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return locations, tx.Commit()
}

func getLocationsByUserID(ctx context.Context, tx *Tx, userId string) ([]*app.Location, error) {
	var defaultLocation sql.NullString
	rows, err := tx.QueryContext(ctx, `
		SELECT
			locations.location_id,
			locations.name,
			locations.address,
			locations.latitude,
			locations.longitude,
			users.location_id AS default_location_id
		FROM locations
		LEFT JOIN users ON locations.user_id = users.user_id
		WHERE locations.user_id = ?
		ORDER BY locations.created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	locations := make([]*app.Location, 0)
	for rows.Next() {
		var location app.Location
		if err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.Address,
			&location.Latitude,
			&location.Longitude,
			&defaultLocation,
		); err != nil {
			log.Println("rows scan error:", err)
			return nil, err
		}
		// mark default location
		if location.ID == defaultLocation.String {
			location.Default = true
		} else {
			location.Default = false
		}
		locations = append(locations, &location)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return locations, nil
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
	_, err := tx.ExecContext(ctx, `
	INSERT INTO locations(
		location_id,
		user_id,
		name,
		address,
		latitude,
		longitude
		) VALUES (?,?,?,?,?,?)
		`,
		location.ID,
		location.UserID,
		location.Name,
		location.Address,
		location.Latitude,
		location.Longitude,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocationService) RemoveLocation(ctx context.Context, locationID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := removeLocation(ctx, tx, locationID); err != nil {
		return err
	}
	return tx.Commit()
}

func removeLocation(ctx context.Context, tx *Tx, locationID string) error {
	_, err := tx.ExecContext(ctx, `
	DELETE FROM locations WHERE location_id = ?
	`, locationID)
	if err != nil {
		return err
	}
	return nil
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
			&location.Name,
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

func findBookingByID(ctx context.Context, tx *Tx, id uuid.UUID) (_ *app.Booking, err error) {
	booking := &app.Booking{}
	// Execue query with limiting WHERE clause and LIMIT/OFFSET injected.
	if err := tx.QueryRowContext(ctx, `
		SELECT
			bookings.booking_id,
			bookings.status,
			bookings.start_at,
			bookings.client_id,
			bookings.provider_id,
			bookings.created_at,
			services.name
		FROM
			bookings
		LEFT JOIN services ON services.id = bookings.service_id
		WHERE
			bookings.booking_id = ?
		ORDER BY
			bookings.start_at ASC
		`,
		id,
	).Scan(
		&booking.ID,
		&booking.Status,
		&booking.StartAt,
		&booking.Client.UserID,
		&booking.Provider.ID,
		&booking.BookedAt,
		&booking.Service.Name,
	); err != nil {
		return nil, err
	}

	booking.Title = booking.Service.Name

	return booking, nil
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

func findBookings(ctx context.Context, tx *Tx) ([]*app.BookingBrief, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    bookings.booking_id,
			bookings.status,
			bookings.created_at,
			bookings.start_at,
			users.first_name,
			users.last_name,
			locations.name
		FROM bookings
		LEFT JOIN providers ON bookings.provider_id = providers.provider_id
		LEFT JOIN users ON providers.user_id = users.user_id
		LEFT JOIN locations ON bookings.location_id = locations.location_id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var firstname sql.NullString
	var lastname sql.NullString

	bookings := make([]*app.BookingBrief, 0)
	for rows.Next() {
		var booking app.BookingBrief
		if err := rows.Scan(
			&booking.ID,
			&booking.Status,
			&booking.BookedAt,
			&booking.StartAt,
			&firstname,
			&lastname,
			&booking.Location,
		); err != nil {
			return nil, err
		}
		booking.Provider = strings.TrimSpace(firstname.String + " " + lastname.String)
		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
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

// createBooking creates a new booking.
func createBooking(ctx context.Context, tx *Tx, booking *model.Booking) error {
	booking.Status = statusPending

	query := `
	INSERT INTO bookings (
		booking_id,
		status,
		start_at,
		client_id,
		provider_id,
		location_id,
		service_id
	) VALUES (?,?,?,?,?,?,?)
	`

	// Insert row into database.
	_, err := tx.ExecContext(ctx, query,
		booking.ID,
		booking.Status,
		booking.StartDate,
		booking.ClientID,
		booking.ProviderID,
		booking.LocationID,
		booking.ServiceID,
	)
	if err != nil {
		return err
	}
	return nil
}

func createPhoto(ctx context.Context, tx *Tx, photo model.Photo) error {
	photo.ID = uuid.New()
	_, err := tx.ExecContext(ctx, `INSERT INTO photos(
		photo_id,
		uploaded_by,
		photo_url,
		booking_id,
		portfolio_id
		) VALUES (?,?,?,?,?)
		`,
		photo.ID,
		photo.OwnerID,
		photo.Url,
		photo.BookingID,
		photo.PortfolioID,
	)
	if err != nil {
		return err
	}
	return nil
}

type BidService struct {
	db *DB
}

func NewBidService(db *DB) *BidService {
	return &BidService{db}
}

func (s *BidService) CreateBid(ctx context.Context, bid *model.Bid) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createBid(ctx, tx, bid); err != nil {
		return err
	}
	return tx.Commit()
}

func createBid(ctx context.Context, tx *Tx, bid *model.Bid) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO bids (
			booking_id,
			provider_id,
			amount
		) VALUES (?, (SELECT provider_id FROM providers WHERE user_id = ?), ?)
		`,
		bid.BookingID,
		bid.BidderID,
		bid.Amount,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *BidService) FindBidsByBookingID(ctx context.Context, bookingID string) ([]*app.Bid, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bids, err := listBids(ctx, tx, bookingID)
	if err != nil {
		return nil, err
	}
	return bids, tx.Commit()
}

func listBids(ctx context.Context, tx *Tx, bookingID string) ([]*app.Bid, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT 
			bids.id,
			bids.booking_id,
			bids.provider_id,
			bids.amount,
			bids.updated_at,
			users.first_name,
			users.last_name
		FROM bids
		LEFT JOIN providers ON bids.provider_id = providers.provider_id
		LEFT JOIN users ON providers.user_id = users.user_id
		WHERE bids.booking_id = ?
		`,
		bookingID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var firstname sql.NullString
	var lastname sql.NullString

	bids := make([]*app.Bid, 0)
	for rows.Next() {
		var bid app.Bid
		if err := rows.Scan(
			&bid.ID,
			&bid.BookingID,
			&bid.Provider.ID,
			&bid.Amount,
			&bid.Date,
			&firstname,
			&lastname,
		); err != nil {
			return nil, err
		}
		bid.Provider.Name = strings.TrimSpace(firstname.String + " " + lastname.String)
		bids = append(bids, &bid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bids, nil
}

func (s *BidService) ListMyBids(ctx context.Context, userID string) ([]*app.Bid, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bids, err := listBidsByCriteria(ctx, tx, userID, "1", "1")
	if err != nil {
		return nil, err
	}
	return bids, tx.Commit()
}

func (s *BidService) FindBidsByRequestID(ctx context.Context, userID string, requestID string) ([]*app.Bid, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bids, err := listBidsByCriteria(ctx, tx, userID, "booking_id", requestID)
	if err != nil {
		return nil, err
	}
	return bids, tx.Commit()
}

func listBidsByCriteria(ctx context.Context, tx *Tx, userID string, haystack string, needle string) ([]*app.Bid, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		bids.id,
		bids.booking_id,
		bids.amount,
		bids.updated_at,
		bids.provider_id,
		users.first_name,
		users.last_name,
		users.photo_url,
		providers.ratings_average,
		providers.reviews_count
	FROM bids
	LEFT JOIN providers ON bids.provider_id = providers.provider_id
	LEFT JOIN users ON providers.user_id = users.user_id
	WHERE bids.provider_id = (SELECT provider_id FROM providers WHERE user_id = ?)
	AND `+haystack+` = ?
	`,
		userID, needle,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bids := make([]*app.Bid, 0)
	for rows.Next() {
		var firstname sql.NullString
		var lastname sql.NullString
		var bid app.Bid
		if err := rows.Scan(
			&bid.ID,
			&bid.BookingID,
			&bid.Amount,
			&bid.Date,
			&bid.Provider.ID,
			&firstname,
			&lastname,
			&bid.Provider.Photo,
			&bid.Provider.Rating,
			&bid.Provider.Reviews,
		); err != nil {
			return nil, err
		}
		bid.Provider.Name = strings.TrimSpace(firstname.String + " " + lastname.String)
		bids = append(bids, &bid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bids, nil
}

func (s *RequestService) FilterRequests(ctx context.Context, filter model.RequestFilter) ([]app.Request, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	requests, err := filterRequests(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func filterRequests(ctx context.Context, tx *Tx, filter model.RequestFilter) (_ []app.Request, err error) {
	// Build WHERE clause. Each part of the WHERE clause is AND-ed together.
	// Values are appended to an arg list to avoid SQL injection.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ClientID; v != "" {
		where = append(where, "client_id = ?")
		args = append(args, v)
	}
	if v := filter.Status; v != "" {
		where, args = append(where, "status = ?"), append(args, v)
	}

	requests := []app.Request{}
	rows, err := tx.QueryContext(ctx, `
		SELECT
			booking_id,
			title,
			status,
			start_at,
			created_at
		FROM bookings
		WHERE `+strings.Join(where, " AND ")+`
	`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		request := app.Request{}
		if err := rows.Scan(
			&request.ID,
			&request.Title,
			&request.Status,
			&request.Start,
			&request.Created,
		); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

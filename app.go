package app

import (
	"context"

	"github.com/andrwkng/hudumaapp/model"
	"github.com/google/uuid"
)

var UserIDNil UserID

type UserID string

func (id UserID) String() string {
	return string(id)
}

type Config struct {
	DB struct {
		DSN string `toml:"dsn"`
	}
}

type App struct {
	// SQLite database used by SQLite service implementations.
	//DB *sqlite.DB
	// Configuration path and parsed config data.
	//config Config
	// HTTP server for handling HTTP communication.
	// SQLite services are attached to it before running.
	//HTTPServer *http.Server
}

func NewApp() *App {
	return &App{
		//Config:     DefaultConfig(),
		//ConfigPath: DefaultConfigPath,

		//DB:         sqlite.NewDB(""),
		//HTTPServer: http.NewServer(),
	}
}

/*
func (a *App) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if m.DB != nil {
		if err := a.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}
*/

type AuthUser struct {
	ID   string  `json:"user_id"`
	Name *string `json:"name"`
	//PictureUrl    *string `json:"picture_url"`
	Provider string `json:"provider"`
	//Email         *string `json:"email_address"`
	EmailVerified bool   `json:"email_verified"`
	PhoneNumber   string `json:"phone_number"`
}

type BookingService interface {
	FindBookingByID(context.Context, uuid.UUID) (*Booking, error)
	CreateBooking(context.Context, *model.Booking) error
	FindBookings(context.Context) ([]*BookingBrief, error)
}

type CategoryService interface {
	CreateCategory(context.Context, *model.Category) error
	ListCategories(context.Context) ([]*Category, error)
}

type UserService interface {
	CreateUser(context.Context, *model.User) error
	FindUserByID(context.Context, string) (*User, error)
	FindUserByUsername(context.Context, string) (*User, error)
	FindUserByPhoneNumber(context.Context, string) (*User, error)
	ValidateUser(context.Context, string, string) error
	UpdateUserPassword(context.Context, *model.ResetUser) error
	UpdateResetCode(context.Context, int, string) error
	UpdateUser(context.Context, *model.User) error
	// Provider
	CreateProvider(context.Context, *model.Provider) error
	FindProviderByID(context.Context, string) (*Provider, error)
	FindProviderByUserID(context.Context, string) (*Provider, error)
	ListProviders(context.Context) ([]*ProviderBrief, error)
	UpdateProvider(context.Context, *model.Provider) error
	// User profile
	FindProfileByUserID(context.Context, string) (*Profile, error)
	CreateProfile(context.Context, *model.Profile) error
	UpdateProfile(context.Context, *model.Profile) error
}

type ClientService interface {
	FindClientByID(context.Context, uuid.UUID) (*Client, error)
	FindClients(context.Context) ([]*Client, error)
}

type ServiceService interface {
	CreateService(context.Context, *model.Service) error
	//FindServiceByID(context.Context, uuid.UUID) (*Service, error)
	ListMyServices(context.Context, string) ([]*Service, error)
	ListServices(context.Context) ([]*Service, error)
}

type ReviewService interface {
	FindReviews(context.Context) ([]*Review, error)
	CreateReview(context.Context, *model.Review) error
}

type PortfolioService interface {
	CreatePortfolio(context.Context, *model.Portfolio) error
	FindPortfolioByID(context.Context, uuid.UUID) (*Portfolio, error)
	ListPortfoliosByUserId(context.Context, string) ([]*Portfolio, error)
}

type RequestService interface {
	FindRequestByID(context.Context, uuid.UUID) (*RequestDetail, error)
	ListRequests(context.Context, UserID) ([]*Request, error)
	CreateRequest(context.Context, *model.Request) error
}

type BidService interface {
	FindBids(context.Context) ([]*Bid, error)
	CreateBid(context.Context, *model.Bid) error
}

type LocationService interface {
	CreateLocation(context.Context, *model.Location) error
	//FindLocationByID(context.Context, uuid.UUID) (*Location, error)
	ListMyLocations(context.Context, string) ([]*Location, error)
}

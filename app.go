package app

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/google/uuid"
	"google.golang.org/api/option"
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
	Fbase *firebase.App
}

func NewApp() *App {
	return &App{
		//Config:     DefaultConfig(),
		//ConfigPath: DefaultConfigPath,

		//DB:         sqlite.NewDB(""),
		//HTTPServer: http.NewServer(),
		Fbase: newFirebaseApp(context.Background()),
	}
}

func newFirebaseApp(ctx context.Context) *firebase.App {
	opt := option.WithCredentialsFile("keys/hudumaapp-firebase-adminsdk-jtet8-7370576c3f.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("error creating firebase app: %s\n", err)
		return nil
	}
	return app
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
	FindProviderBookingByID(context.Context, uuid.UUID, string) (*ProviderBooking, error)
	CreateBooking(context.Context, *model.Booking) error
	FindMyBookings(context.Context) ([]*BookingBrief, error)
	FindBookings(context.Context, string) ([]*BookingBrief, error)
	InsertDate(context.Context, string) error
}

type CategoryService interface {
	CreateCategory(context.Context, *model.Category) error
	ListCategories(context.Context) ([]*Category, error)
	ListRootCategories(context.Context) ([]RootCategory, error)
	ListCategoriesByParentID(context.Context, string) ([]*Category, error)
	ListCategoriesByIndustryID(context.Context, string) ([]*Category, error)
}

type IndustryService interface {
	CreateIndustry(context.Context, *model.Industry) error
	ListIndustries(context.Context) ([]*Industry, error)
}

type UserService interface {
	// User
	CreateUser(context.Context, *model.User) error
	FindUserByID(context.Context, string) (*User, error)
	FindUserByUsername(context.Context, string) (*User, error)
	FindUserByPhoneNumber(context.Context, string) (*User, error)
	// Password
	ValidateUser(context.Context, string, string, bool) error
	ResetUserPassword(context.Context, string, string) error
	ChangeUserPassword(context.Context, *model.PwdChange) error
	// Provider
	CreateProvider(context.Context, *model.Provider) error
	FindProviderByID(context.Context, string) (*Provider, error)
	FindProviderByUserID(context.Context, string) (*Provider, error)
	ListProviders(context.Context) ([]*ProviderBrief, error)
	FilterProviders(context.Context, model.ProviderFilter) ([]*ProviderBrief, error)
	UpdateProvider(context.Context, *model.Provider) error
	// User profile
	FindProfileByUserID(context.Context, string) (*Profile, error)
	CreateProfile(context.Context, *model.Profile) error
	UpdateProfile(context.Context, *model.Profile) error
	// Service
	CreateService(context.Context, *model.Service) error
	//FindServiceByID(context.Context, string) (*Service, error)
	ListMyServices(context.Context, string) ([]*Service, error)
	ListServicesByProviderID(context.Context, string) ([]*Service, error)
}

type ClientService interface {
	FindClientByID(context.Context, uuid.UUID) (*Client, error)
	FindClients(context.Context) ([]*Client, error)
}

type ReviewService interface {
	//FindReviews(context.Context) ([]*Review, error)
	CreateReview(context.Context, *model.Review) error
	ListReviewsByProviderID(context.Context, string) ([]*Review, error)
}

type PortfolioService interface {
	CreatePortfolio(context.Context, *model.Portfolio) error
	FindPortfolioByID(context.Context, uuid.UUID) (*Portfolio, error)
	ListPortfoliosByProviderId(context.Context, string) ([]*PortfolioBrief, error)
	ListPortfoliosByUserId(context.Context, string) ([]*PortfolioBrief, error)
}

type RequestService interface {
	FindRequestByID(context.Context, uuid.UUID) (*RequestDetail, error)
	CreateRequest(context.Context, *model.Request) error
	FilterRequests(context.Context, model.RequestFilter) ([]Request, error)
	AllRequests(context.Context, model.RequestFilter) ([]AllRequest, error)
	ListRequestsCategories(context.Context) ([]Category, error)
}

type BidService interface {
	ListMyBids(context.Context, string) ([]*Bid, error)
	FindBidsByBookingID(context.Context, string) ([]*Bid, error)
	FindBidsByRequestID(context.Context, string, string) ([]*Bid, error)
	CreateBid(context.Context, *model.Bid) error
	AcceptBid(context.Context, int) error
}

type LocationService interface {
	CreateLocation(context.Context, *model.Location) error
	//FindLocationByID(context.Context, uuid.UUID) (*Location, error)
	ListMyLocations(context.Context, string) ([]*Location, error)
	RemoveLocation(context.Context, string) error
}

type SearchService interface {
	SearchByQuery(context.Context, model.Search) ([]SearchResult, error)
	InstantSearchRequests(context.Context, model.Search) ([]RequestSearchResult, error)
}

package app

import (
	"context"
	"time"

	"github.com/andrwkng/hudumaapp/model"
	"github.com/google/uuid"
)

type Appointment struct {
	StartTime time.Time `json:"start_time"`
}

type Rate struct {
	Price    *string `json:"price"`
	Currency *string `json:"currency"`
	Unit     *string `json:"unit"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Provider struct {
	ID int `json:"id"`
	Profile
	Bio        *string `json:"bio"`
	Profession *string `json:"profession"`
	//Professions []string `json:"professons"`
	AvgRating float32 `json:"rating"`
	Stats     Stats   `json:"stats"`
	//distanceKM int `json:"distance`
	Rate     `json:"rate"`
	Services []Service `json:"services"`
}

type ProviderBrief struct {
	ID       uuid.UUID `json"id"`
	UserID   string    `json:"user_id"`
	Name     string    `json:"name"`
	Rate     `json:"rate"`
	Jobs     int     `json:"num_jobs"`
	Rating   float32 `json:"avg_rating"`
	Reviews  int     `json:"num_reviews"`
	Photo    string  `json:"photo_url"`
	Distance string  `json:"distance"`
}

type Stats struct {
	Reviews    int `json:"reviews"`
	Jobs       int `json:"jobs"`
	Portfolios int `json:"portfolios"`
	Services   int `json:"services"`
}

type Location struct {
	ID        *string  `json:"location_id"`
	Name      *string `json:"name"`
	Title     *string `json:"title"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

type Service struct {
	//ID    uuid.UUID `json:"service_id"
	Title *string `json:"service_title"`
}

type Review struct {
	ID        uuid.UUID `json:"review_id"`
	Provider  User      `json:"provider"`
	Client    User      `json:"client"`
	Service   Service   `json:"service"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type Request struct {
	ID      uuid.UUID `json:"request_id"`
	Title   string
	Status  string
	Created string
	Start   string
	Bids    int
}

type RequestDetail struct {
	ID      uuid.UUID `json:"request_id"`
	Title   string    `json:"title"`
	Note    string    `json:"note"`
	Status  string    `json:"status"`
	Created string    `json:"created"`
	Start   string    `json:"start"`
	Bids    int       `json:"bids"`
	Photos  []string  `json:"photos"`
	Client  string    `json:"client"`
}

type Client struct {
	Profile
}

type Booking struct {
	ID          uuid.UUID `json:"booking_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	Description *string   `json:"description"`
	Type        *string   `json:"type"`
	BookedAt    string    `json:"booked_at"`
	Photos      []string  `json:"photos"`
	StartAt     string    `json:"start_time"`
	Category    string    `json:"category"`
	Service
	Provider `json:"provider"`
	Client   `json:"client"`
	Location `json:"location"`
}

type BookingBrief struct {
	ID          uuid.UUID `json:"booking_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	Description string    `json:"descripton"`
	Type        string    `json:"type"`
	BookedAt    string    `json:"bookedat"`
	Photos      []string  `json:"photos"`
	StartAt     string    `json:"start_tie"`
	Category    string    `json:"category"`
	Service     string    `json:"service"`
	Provider    string    `json:"proider"`
	Client      string    `json:"clent"`
	Location    string    `json:"loction"`
}

type Portfolio struct {
	ID          uuid.UUID `json:"portfolio_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	Description string    `json:"descripton"`
	Type        string    `json:"type"`
	Photos      []string  `json:"photos"`
	StartAt     string    `json:"start_tie"`
	Category    string    `json:"category"`
	Service     string    `json:"service"`
	Provider    string    `json:"provider"`
	BookingID   string    `json:"booking_id"`
}

type User struct {
	UserID      string  `json:"user_id"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	DisplayName *string `json:"display_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Location    `json:"location"`
	PhotoUrl    *string `json:"photo_url"`
}

type Profile struct {
	ID uuid.UUID `json:"profile_id"`
	User
	//UserID        string  `json:"user_id"`
	//FirstName     *string `json:"first_name"`
	//LastName      *string `json:"last_name"`
	//DisplayName   *string `json:"display_name"`
	//Email         *string `json:"email"`
	EmailVerified bool `json:"email_verified"`
	//Phone         *string `json:"phone"`
	Verified  bool    `json:"verified"`
}

type Bid struct{}

type Transaction struct {
	Code      string
	Amount    int
	Currency  string
	Status    string
	CreatedAt time.Time
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
	// Provider
	CreateProvider(context.Context, *model.Provider) error
	FindProviderByID(context.Context, string) (*Provider, error)
	ListProviders(context.Context) ([]*ProviderBrief, error)
	// User profile
	GetProfile(context.Context, string) (*Profile, error)
	FindProfileByUserID(context.Context, string) (*Profile, error)
	CreateProfile(context.Context, *model.Profile) error
	UpdateProfile(context.Context, *model.Profile) error
}

type ClientService interface {
	FindClientByID(context.Context, uuid.UUID) (*Client, error)
	FindClients(context.Context) ([]*Client, error)
}

type ServiceService interface {
	FindServiceByID(context.Context, uuid.UUID) (*Service, error)
	FindServices(context.Context) ([]*Service, error)
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
	ListMyLocations(context.Context, *AuthUser) ([]*Location, error)
}

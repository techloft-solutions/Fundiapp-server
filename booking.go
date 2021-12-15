package app

import (
	"time"

	"github.com/google/uuid"
)

type Price struct {
	Amount   *int    `json:"amount"`
	Currency *string `json:"currency"`
}

type Category struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
	IconURL  string  `json:"icon_url"`
}

type RootCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Industry struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}

type Provider struct {
	ID string `json:"provder_id"`
	Profile
	Bio        *string `json:"bio"`
	Profession *string `json:"profession"`
	//Professions []string `json:"professons"`
	AvgRating float32    `json:"rating"`
	Stats     Stats      `json:"stats"`
	Price     *Price     `json:"price"`
	Services  []*Service `json:"services"`
	Phone     string     `json:"phone"`
}

type ProviderBrief struct {
	ID         uuid.UUID `json:"provder_id"`
	Name       string    `json:"name"`
	Profession *string   `json:"profession"`
	Price      *Price    `json:"price"`
	Jobs       int       `json:"num_jobs"`
	Rating     float32   `json:"avg_rating"`
	Reviews    int       `json:"num_reviews"`
	Photo      *string   `json:"photo_url"`
}

type SearchResult struct {
	ID   uuid.UUID `json:"provder_id"`
	Name string    `json:"name"`
	//Rating   float32   `json:"avg_rating"`
	//Reviews  int       `json:"num_reviews"`
	Photo    *string `json:"photo_url"`
	Distance *string `json:"distance"`
}

type RequestSearchResult struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"name"`
	Count        int    `json:"count"`
}

type Stats struct {
	Reviews    int `json:"reviews"`
	Jobs       int `json:"jobs"`
	Portfolios int `json:"portfolios"`
	Services   int `json:"services"`
}

type Location struct {
	ID        string  `json:"location_id"`
	Name      *string `json:"name"`
	Latitude  string  `json:"latitude"`
	Longitude string  `json:"longitude"`
	Address   *string `json:"address"`
	Default   bool    `json:"default"`
}

type Service struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    `json:"price"`
	Category *string `json:"category"`
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

type RequestProvider struct {
	ID    *uuid.UUID `json:"id"`
	Name  *string    `json:"name"`
	Photo *string    `json:"photo_url"`
}

type Request struct {
	ID        uuid.UUID `json:"request_id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created"`
	StartAt   string    `json:"start"`
	Bids      int       `json:"bids"`
}

type AllRequest struct {
	ID        uuid.UUID `json:"request_id"`
	Title     string    `json:"title"`
	Category  *string   `json:"category"`
	Urgent    bool      `json:"urgent"`
	Distance  string    `json:"distance_km"`
	CreatedAt string    `json:"created_at"`
	StartAt   string    `json:"start_at"`
	Address   string    `json:"location"`
}

type location struct {
	ID        string `json:"location_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Address   string `json:"address"`
}

type RequestDetail struct {
	ID        uuid.UUID `json:"request_id"`
	Title     string    `json:"title"`
	Category  *string   `json:"category"`
	Note      string    `json:"note"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"posted"`
	StartAt   string    `json:"start_at"`
	Bids      int       `json:"bids"`
	Photos    []string  `json:"photos"`
	Location  location  `json:"location"`
}

type Client struct {
	Profile
}

type bookingLocation struct {
	ID        string `json:"location_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Address   string `json:"address"`
}

type bookingService struct {
}

type bookingUser struct {
	UserID    string  `json:"user_id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Username  *string `json:"display_name"`
	Phone     *string `json:"phone"`
	PhotoUrl  *string `json:"photo_url"`
}

type bookingProvider struct {
	ProviderID string `json:"provider_id"`
	bookingUser
}

type Booking struct {
	ID          uuid.UUID `json:"booking_id"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	Description *string   `json:"description"`
	//Type     *string  `json:"type"`
	BookedAt string `json:"booked_at"`
	//Photos   []string `json:"photos"`
	StartAt string `json:"start_at"`
	//Service  `json:"service"`
	Provider bookingProvider `json:"provider"`
	Client   bookingUser     `json:"client"`
	//Location bookingLocation `json:"location"`
}

type BookingBrief struct {
	ID     uuid.UUID `json:"booking_id"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
	//Description string    `json:"descripton"`
	//Type     string `json:"type"`
	BookedAt string `json:"booked_at"`
	//Photos      []string  `json:"photos"`
	StartAt string `json:"start_at"`
	//Category string `json:"category"`
	//Service  string `json:"service"`
}

type Portfolio struct {
	ID         uuid.UUID `json:"portfolio_id"`
	Title      string    `json:"title"`
	Photos     []string  `json:"photos"`
	Service    string    `json:"service,omitempty"`
	ProviderID string    `json:"provider,omitempty"`
	//BookingID  string    `json:"booking_id"`
}

type PortfolioBrief struct {
	ID    uuid.UUID `json:"portfolio_id"`
	Title string    `json:"title"`
}

type User struct {
	UserID     string  `json:"user_id"`
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Username   *string `json:"display_name"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	PhotoUrl   *string `json:"photo_url"`
	IsProvider bool    `json:"-"`
}

type ProfileLocation struct {
	ID      *string `json:"location_id"`
	Address *string `json:"location_address"`
}

type Profile struct {
	User
	// UserID    string  `json:"user_id"`
	// FirstName *string `json:"first_name"`
	// LastName  *string `json:"last_name"`
	// Username  *string `json:"display_name"`
	// Email     *string `json:"email"`
	// Phone     *string `json:"phone"`
	Location      *ProfileLocation `json:"location"`
	PhotoUrl      *string          `json:"photo_url"`
	EmailVerified bool             `json:"email_verified"`
	Verified      bool             `json:"verified"`
}

type Bid struct {
	ID        int           `json:"bid_id"`
	BookingID uuid.UUID     `json:"request_id"`
	Provider  ProviderBrief `json:"bidder"`
	Amount    int           `json:"amount"`
	Date      string        `json:"date"`
}

type BidBrief struct {
	ID         int       `json:"bid_id"`
	BookingID  uuid.UUID `json:"request_id"`
	ProviderID string    `json:"provider_id"`
	Bidder     string    `json:"bidder"`
	Amount     int       `json:"amount"`
	Date       string    `json:"date"`
}

type Transaction struct {
	Code      string
	Amount    int
	Currency  string
	Status    string
	CreatedAt time.Time
}

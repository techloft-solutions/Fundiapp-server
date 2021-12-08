package app

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	StartTime time.Time `json:"start_time"`
}

type price struct {
	Amount   int
	Currency string
}

type Rate struct {
	Price    *string `json:"price"`
	Currency *string `json:"currency"`
	Unit     *string `json:"unit"`
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
	AvgRating float32 `json:"rating"`
	Stats     Stats   `json:"stats"`
	Rate      `json:"rate"`
	Services  []*Service `json:"services"`
}

type ProviderBrief struct {
	ID         uuid.UUID `json:"provder_id"`
	Name       string    `json:"name"`
	Profession *string   `json:"profession"`
	Rate       `json:"rate"`
	Jobs       int     `json:"num_jobs"`
	Rating     float32 `json:"avg_rating"`
	Reviews    int     `json:"num_reviews"`
	Photo      *string `json:"photo_url"`
}

type SearchResult struct {
	ID   uuid.UUID `json:"provder_id"`
	Name string    `json:"name"`
	//Rating   float32   `json:"avg_rating"`
	//Reviews  int       `json:"num_reviews"`
	Photo    *string `json:"photo_url"`
	Distance *string `json:"distance"`
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
	Rate     `json:"rate"`
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
	ID        uuid.UUID        `json:"request_id"`
	Title     string           `json:"title"`
	Status    string           `json:"status"`
	CreatedAt string           `json:"created"`
	StartAt   string           `json:"start"`
	Bids      int              `json:"bids"`
	Provider  *RequestProvider `json:"provider"`
}

type location struct {
	ID        string `json:"location_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Address   string `json:"address"`
}

type RequestDetail struct {
	ID       uuid.UUID `json:"request_id"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Note     string    `json:"note"`
	Status   string    `json:"status"`
	Created  string    `json:"created"`
	Start    string    `json:"start"`
	Bids     int       `json:"bids"`
	Photos   []string  `json:"photos"`
	Location location  `json:"location"`
	Client   string    `json:"client"`
}

type Client struct {
	Profile
}

type Booking struct {
	ID     uuid.UUID `json:"booking_id"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
	//Description *string  `json:"description"`
	//Type     *string  `json:"type"`
	BookedAt string `json:"booked_at"`
	//Photos   []string `json:"photos"`
	StartAt  string `json:"start_time"`
	Category string `json:"category"`
	Service  `json:"service"`
	Provider ProviderBrief `json:"provider"`
	Client   `json:"client"`
	Location `json:"location"`
}

type BookingBrief struct {
	ID     uuid.UUID `json:"booking_id"`
	Status string    `json:"status"`
	//Description string    `json:"descripton"`
	//Type     string `json:"type"`
	BookedAt string `json:"booked_at"`
	//Photos      []string  `json:"photos"`
	StartAt string `json:"start_time"`
	//Category string `json:"category"`
	//Service  string `json:"service"`
	Provider string `json:"provider"`
	Location string `json:"location"`
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

type Profile struct {
	User
	// UserID    string  `json:"user_id"`
	// FirstName *string `json:"first_name"`
	// LastName  *string `json:"last_name"`
	// Username  *string `json:"display_name"`
	// Email     *string `json:"email"`
	// Phone     *string `json:"phone"`
	LocationID    *string `json:"location_id"`
	Address       *string `json:"location_address"`
	PhotoUrl      *string `json:"photo_url"`
	EmailVerified bool    `json:"email_verified"`
	Verified      bool    `json:"verified"`
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

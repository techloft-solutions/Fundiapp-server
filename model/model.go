package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID `valid:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type User struct {
	Username string `valid:"required" json:"displayname,omitempty"`
	Password string `valid:"required" json:"password"`
	Phone    string `valid:"required" json:"phone"`
}

type Provider struct {
	ID uuid.UUID `valid:"required" json:"provider_id"`
	Profile
	Bio        *string `json:"bio"`
	Profession *string `json:"profession"`
}

type Service struct {
	UserID     string
	ProviderID uuid.UUID
	Name       string `json:"name"`
	Rate
}

type Rate struct {
	Amount   uint   `json:"rate_amount"`
	Unit     string `json:"rate_unit"`
	Currency string `json:"rate_currency"`
}

type Statistic struct {
	Model
	ProviderID     uuid.UUID
	RatingsAvg     int
	ReviewsCount   int
	ServicesCount  int
	PortfolioCount int
}

type Review struct {
	AuthorID             string  `valid:"required" json:"author_id"`
	ProviderID           string  `valid:"required,uuid" json:"provider_id"`
	ServiceID            string  `valid:"required,uuid" json:"service_id"`
	Comment              string  `valid:"required" json:"comment"`
	Rating               string  `valid:"required" json:"rating"`
	IntegrityRating      *string `json:"integrity_rating"`
	CompetenceRating     *string `json:"competence_rating"`
	ResponsivenessRating *string `json:"responsiveness_rating"`
	QualityRating        *string `json:"quality_rating"`
}

type Booking struct {
	Model
	Title       string `valid:"required"`
	StartDate   string `valid:"required,rfc3339"`
	Description *string
	LocationID  string `valid:"required,uuid"`
	Status      string
	ProviderID  *string
	ClientID    string `valid:"required"`
	ServiceID   string `valid:"required,uuid"`
	Photos      []string
}

type Request struct {
	ID         uuid.UUID `valid:"required"`
	Title      string    `valid:"required" json:"title"`
	StartDate  string    `valid:"required, rfc3339" json:"start_date"`
	Note       string    `valid:"required" json:"note"`
	LocationID string    `valid:"required,uuid" json:"location_id"`
	Type       string    `valid:"required" json:"type"`
	ClientID   string    `valid:"required" json:"client_id"`
	Photos     []string  `json:"photos"`
	Status     string    `json:"status"`
	Urgent     bool      `valid:"required" json:"urgent,string"`
}

type Photo struct {
	Model
	Owner string
	Url   string
}

type BookingPhoto struct {
	BookingID uuid.UUID
	PhotoID   uuid.UUID
}

type PortfolioPhoto struct {
	ID string
}

type Portfolio struct {
	Model
	Title     string
	BookingID uuid.UUID
	Photos    []string
}

type Category struct {
	Name        string `valid:"required"`
	Description *string
	ParentID    *int
	Profession  *string
}

type ProviderProfession struct {
	Model
	ProviderID uuid.UUID
	CategoryID uuid.UUID
}

type Location struct {
	ID        uuid.UUID `json:"location_id" valid:"required"`
	Name      *string   `json:"name,omitempty"`
	Latitude  string    `valid:"required" json:"latitude"`
	Longitude string    `valid:"required" json:"longitude"`
	City      *string   `json:"city,omitempty"`
	State     *string   `json:"state,omitempty"`
	Zip       *string   `json:"zip,omitempty"`
	UserID    string    `valid:"required"`
	Address   string    `valid:"required" json:"address"`
}

type BookingLocation struct {
	Model
	BookingID  uuid.UUID
	LocationID uuid.UUID
}

type UserLocation struct {
	Model
	LocationID uuid.UUID
	UserID     string
}

type Bid struct {
	BookingID uuid.UUID
	BidderID  string
	Price     int `valid:"required"`
}

type Transaction struct {
	Model
	Code       string
	ServiceID  uuid.UUID
	BookingID  uuid.UUID
	UserID     string
	ProviderID string
	Amount     int
	Currency   string
}

type Profile struct {
	//Model
	ID     uuid.UUID `valid:"required" json:"profile_id"`
	UserID string    `valid:"required" json:"user_id"`
	//Username  *string   `json:"display_name,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	//Phone       string    `json:"phone,omitempty"`
	PhotoUrl   *string `json:"photo_url,omitempty"`
	LocationID *string `json:"location_id,omitempty"`
	Status     *string `json:"status,omitempty"`
	Type       string  `valid:"required,omitempty"`
	Verified   bool
}

type Schedule struct {
	StartTime string
	EndTime   string
}

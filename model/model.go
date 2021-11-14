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
	UserID     string
	LocationID uuid.UUID
}

type Provider struct {
	Profile
	Bio        *string `json:"bio"`
	Profession *string `json:"profession"`
}

type Client struct {
	Model
	User
}

type Service struct {
	Model
	User
	ProviderID uuid.UUID
	ClientID   uuid.UUID
}

type Rate struct {
	Model
	Amount   uint
	Unit     string
	Currency string
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
	Model
	AuthorID             uuid.UUID
	ProviderID           uuid.UUID
	ServiceID            uuid.UUID
	RateID               uuid.UUID
	Comment              string
	Rating               float32
	IntegrityRating      float32
	CompetenceRating     float32
	ResponsivenessRating float32
	QualityRating        float32
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
	Model
	Title      string `valid:"required"`
	StartDate  string `valid:"required,rfc3339"`
	Note       string `valid:"required"`
	LocationID string `valid:"required,uuid"`
	Type       string `valid:"required"`
	ClientID   string `valid:"required"`
	Photos     []string
	Status     string
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
	Description string
	ParentID    int
	Profession  string
}

type ProviderProfession struct {
	Model
	ProviderID uuid.UUID
	CategoryID uuid.UUID
}

type Location struct {
	ID        uuid.UUID `valid:"required"`
	Name      *string   `json:"display_name"`
	Latitude  string    `valid:"required" json:"latitude"`
	Longitude string    `valid:"required" json:"longitude"`
	City      *string   `json:"city"`
	State     *string   `json:"state"`
	Zip       *string   `json:"zip"`
	UserID    string    `valid:"required"`
	Address   *string   `json:"address"`
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
	Model
	BookingID uuid.UUID
	UserID    string
	Price     int
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
	ID          uuid.UUID `valid:"required"`
	UserID      string    `valid:"required"`
	DisplayName *string   `json:"display_name"`
	FirstName   *string   `json:"first_name"`
	LastName    *string   `json:"last_name"`
	Email       *string   `json:"email"`
	Phone       string    `json:"phone"`
	PhotoUrl    *string   `json:"photo_url"`
	LocationID  *string   `json:"location_id"`
	Status      *string
	Type        string `valid:"required"`
	Verified    bool
}

type Schedule struct {
	StartTime string
	EndTime   string
}

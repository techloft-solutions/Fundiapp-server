package model

import (
	"database/sql"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID `valid:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

func (m Model) Validate(v interface{}) error {
	_, err := govalidator.ValidateStruct(v)
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	UserID     string `valid:"required" json:"user_id"`
	Username   string `valid:"required" json:"display_name,omitempty"`
	Password   string `valid:"required" json:"password"`
	Phone      string `valid:"required" json:"phone"`
	IsProvider bool   `json:"is_provider,omitempty"`
}

type PwdChange struct {
	UserID      string `valid:"required" json:"user_id"`
	OldPassword string `valid:"required" json:"old_password"`
	NewPassword string `valid:"required" json:"new_password"`
}

type Provider struct {
	ID uuid.UUID `valid:"required" json:"provider_id"`
	Profile
	Bio        *string `json:"bio"`
	Profession *string `json:"profession"`
	Rate
	CategoryID *string `json:"category_id"`
	IndustryID *string `json:"industry_id"`
}

type Service struct {
	ProviderID uuid.UUID
	Name       string `valid:"required" json:"name"`
	Rate
}

type Rate struct {
	Amount   string `json:"rate_amount"`
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
	Comment              *string `valid:"required" json:"comment"`
	Rating               *string `valid:"required" json:"rating"`
	IntegrityRating      *string `json:"integrity_rating"`
	CompetenceRating     *string `json:"competence_rating"`
	ResponsivenessRating *string `json:"responsiveness_rating"`
	QualityRating        *string `json:"quality_rating"`
}

type Booking struct {
	ID         uuid.UUID `json:"booking_id"`
	StartDate  string    `valid:"required,rfc3339" json:"start_date"`
	LocationID string    `valid:"required,uuid" json:"location_id"`
	Status     string
	ProviderID string `valid:"required" json:"provider_id"`
	ClientID   string `valid:"required" json:"client_id"`
	ServiceID  string `valid:"required,uuid" json:"service_id"`
}

type Request struct {
	ID         uuid.UUID `json:"request_id"`
	Title      string    `valid:"required" json:"title"`
	StartDate  string    `valid:"required, rfc3339WithoutZone" json:"start_date"`
	Note       string    `valid:"required" json:"note"`
	LocationID string    `valid:"required,uuid" json:"location_id"`
	ClientID   string    `valid:"required" json:"client_id"`
	CategoryID *string   `json:"category_id"`
	Photos     []string  `json:"-"`
	Status     string    `json:"status"`
	Urgent     bool      `json:"urgent,string"`
}

type Photo struct {
	ID          uuid.UUID `valid:"required"`
	OwnerID     string    `valid:"required"`
	Url         string    `valid:"required"`
	BookingID   string    `valid:"uuid" json:",omitempty"`
	PortfolioID string    `valid:"uuid" json:",omitempty"`
}

type Portfolio struct {
	ID        uuid.UUID
	OwnerID   string
	UserID    string
	Title     string   `valid:"required" json:"title"`
	BookingID *string  `valid:"uuid" json:",omitempty"`
	ServiceID *string  `valid:"uuid" json:",omitempty"`
	Photos    []string `valid:"required" json:"-"`
}

type Category struct {
	Name        string  `json:"name" valid:"required"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
	IconURL     string  `json:"icon_url" valid:"required"`
}

type Industry struct {
	Name        string  `json:"name" valid:"required"`
	Description *string `json:"description"`
	IconURL     string  `json:"icon_url" valid:"required"`
}

type ProviderProfession struct {
	Model
	ProviderID uuid.UUID
	CategoryID uuid.UUID
}

type Location struct {
	ID        uuid.UUID `json:"location_id" valid:"required"`
	Name      *string   `json:"name,omitempty"`
	Latitude  string    `valid:"required,latitude" json:"latitude"`
	Longitude string    `valid:"required,longitude" json:"longitude"`
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
	BookingID string `valid:"required,uuid" json:"request_id"`
	BidderID  string `valid:"required" json:"user_id"`
	Amount    string `valid:"required" json:"amount"`
}

type Search struct {
	Query     string
	Latitude  string `valid:",latitude"`
	Longitude string `valid:",longitude"`
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
	//ID     uuid.UUID `valid:"required" json:"profile_id"`
	UserID string `valid:"required" json:"user_id"`
	//Username  *string   `json:"display_name,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	//Phone       string    `json:"phone,omitempty"`
	PhotoUrl   *string `valid:"url" json:"photo_url,omitempty"`
	LocationID *string `json:"location_id,omitempty"`
	Status     *string `json:"status,omitempty"`
	Verified   bool
}

type Schedule struct {
	StartTime string
	EndTime   string
}

// ProviderFilter represents a filter used on service providers.
type ProviderFilter struct {
	CategoryID string `json:"category_id"`
	IndustryID string `json:"industry_id"`
}

type RequestFilter struct {
	Status   string
	ClientID string
}

package entity

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type PVZ struct {
	ID               uuid.UUID `json:"pvz_id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city_name"`
	UserID           uuid.UUID
}

type User struct {
	ID       uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

type Receptions struct {
	ID       uuid.UUID  `json:"id"`
	DateTime time.Time  `json:"dateTime"`
	PVZID    uuid.UUID  `json:"pvzId"`
	Status   string     `json:"status"`
	Products []Products `json:"products"`
}

type Products struct {
	ID          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionId uuid.UUID `json:"receptionId"`
}

type ListPVZ struct {
	Pvz        PVZ          `json:"pvz"`
	Receptions []Receptions `json:"receptions"`
}

type Filter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID               uuid.UUID `json:"pvz_id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city_name"`
}

type User struct {
	ID       uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

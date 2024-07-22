package model

import (
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string `gorm:"index"`
	Country   dto.Country
	CreatedAt time.Time
	UpdatedAt time.Time
}

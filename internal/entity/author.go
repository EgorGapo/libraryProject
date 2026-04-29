package entity

import (
	"time"
)

type Author struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

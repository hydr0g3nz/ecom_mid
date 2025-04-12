package common

import (
	"time"
)

// Entity represents a base entity with common fields
type Entity struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

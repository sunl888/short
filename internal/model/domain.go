package model

import "time"

// short domain model.
type Domain struct {
	Id        int       `json:"id"`
	FromUrl   string    `json:"origin"`
	ToUrl     string    `json:"to_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

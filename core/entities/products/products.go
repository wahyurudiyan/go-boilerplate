package entities

import (
	"time"
)

type Product struct {
	Id          uint64 `gorm:"primary_key" json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Price       int64  `json:"price,omitempty"` // store price with the smallest unit like cents
	Description string `json:"description,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type Location struct {
	Id   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type ProductStock struct {
	Id         uint64 `json:"id,omitempty"`
	ProductId  uint64 `json:"product_id,omitempty"`
	LocationId uint64 `json:"location_id,omitempty"`

	Product  Product  `json:"product,omitempty"`
	Location Location `json:"location,omitempty"`
	Qty      uint64   `json:"qty,omitempty"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

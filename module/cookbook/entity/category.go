package entity

import (
	"github.com/guregu/null"
	"time"
)

// Categories is the plural form of Category
type Categories []*Category

// Category holds our category entity
type Category struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt null.Time
	UpdatedBy null.String
	IsDeleted bool
}

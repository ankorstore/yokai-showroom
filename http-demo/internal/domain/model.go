package domain

import (
	"database/sql"
)

// Gopher is the model for gophers.
type Gopher struct {
	ID   int32          `json:"id"`
	Name string         `json:"name"`
	Job  sql.NullString `json:"job"`
}

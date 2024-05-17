package model

import (
	"database/sql"
)

type Gopher struct {
	ID   int32          `json:"id"`
	Name string         `json:"name"`
	Job  sql.NullString `json:"job"`
}

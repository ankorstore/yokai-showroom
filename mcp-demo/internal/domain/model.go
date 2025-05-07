package domain

// Gopher is the model for gophers.
type Gopher struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Job  string `json:"job"`
}

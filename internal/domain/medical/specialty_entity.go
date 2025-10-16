package medical

import "fmt"

// Specialty represents a medical specialty
type Specialty struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (s Specialty) GetId() string {
	return fmt.Sprintf("%d", s.ID)
}

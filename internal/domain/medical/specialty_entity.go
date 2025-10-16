package medical

// Specialty represents a medical specialty
type Specialty struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (s Specialty) GetId() string {
	return s.ID
}

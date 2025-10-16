package medical

// Specialty represents a medical specialty
type Specialty struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

const (
	SpecialtyCardiology       = "cardiology"
	SpecialtyDermatology      = "dermatology"
	SpecialtyEndocrinology    = "endocrinology"
	SpecialtyGastroenterology = "gastroenterology"
	SpecialtyNeurology        = "neurology"
	SpecialtyOncology         = "oncology"
	SpecialtyOrthopedics      = "orthopedics"
	SpecialtyPediatrics       = "pediatrics"
	SpecialtyPsychiatry       = "psychiatry"
	SpecialtyUrology          = "urology"
)

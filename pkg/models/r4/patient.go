package r4

import (
	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
)

// Patient extends the base Patient model with R4-specific fields and validations
type Patient struct {
	models.Patient
}

// Validate performs R4-specific validation rules
func (p *Patient) Validate() error {
	// Add R4-specific validation rules here
	// For example:
	// - Check if required fields are present
	// - Validate field formats
	// - Check value sets
	return nil
}

// ToR5 converts an R4 Patient to R5 format
func (p *Patient) ToR5() (*models.Patient, error) {
	// For now, we'll just return the base patient since we haven't
	// implemented any R4-specific fields yet
	return &p.Patient, nil
}

// FromR5 converts an R5 Patient to R4 format
func (p *Patient) FromR5(r5Patient *models.Patient) error {
	// For now, we'll just copy the base patient since we haven't
	// implemented any R4-specific fields yet
	p.Patient = *r5Patient
	return nil
}

package r4

import (
	"encoding/json"
	"testing"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
)

func TestR4PatientMappingAndValidation(t *testing.T) {
	jsonData := `{
		"resourceType": "Patient",
		"id": "r4-123",
		"active": true,
		"name": [
			{
				"family": "Smith",
				"given": ["Alice"]
			}
		],
		"gender": "female",
		"birthDate": "1985-05-15"
	}`

	var basePatient models.Patient
	err := json.Unmarshal([]byte(jsonData), &basePatient)
	if err != nil {
		t.Fatalf("Failed to unmarshal base patient: %v", err)
	}

	r4Patient := &Patient{Patient: basePatient}

	if err := r4Patient.Validate(); err != nil {
		t.Errorf("R4 Patient validation failed: %v", err)
	}

	if r4Patient.ID != "r4-123" {
		t.Errorf("Expected ID r4-123, got %s", r4Patient.ID)
	}
	if r4Patient.Gender != "female" {
		t.Errorf("Expected gender female, got %s", r4Patient.Gender)
	}
	if len(r4Patient.Name) != 1 || r4Patient.Name[0].Family != "Smith" {
		t.Errorf("Expected family name Smith, got %+v", r4Patient.Name)
	}
}

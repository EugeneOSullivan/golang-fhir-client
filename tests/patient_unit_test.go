package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
)

func TestPatient_JSONRoundTrip(t *testing.T) {
	birth := time.Date(1980, 7, 15, 0, 0, 0, 0, time.UTC)
	patient := models.Patient{
		Base: models.Base{
			ResourceType: models.ResourceTypePatient,
			ID:           "unit-1",
		},
		Active:    ptrBool(true),
		Gender:    "male",
		BirthDate: &birth,
		Name: []models.HumanName{{
			Family: "UnitTest",
			Given:  []string{"Bob"},
		}},
	}
	data, err := json.Marshal(patient)
	if err != nil {
		t.Fatalf("Failed to marshal patient: %v", err)
	}

	var out models.Patient
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("Failed to unmarshal patient: %v", err)
	}
	if out.ID != patient.ID {
		t.Errorf("ID mismatch: got %s, want %s", out.ID, patient.ID)
	}
	if out.Gender != patient.Gender {
		t.Errorf("Gender mismatch: got %s, want %s", out.Gender, patient.Gender)
	}
	if out.BirthDate == nil || !out.BirthDate.Equal(birth) {
		t.Errorf("BirthDate mismatch: got %v, want %v", out.BirthDate, birth)
	}
}

func TestPatient_RequiredFields(t *testing.T) {
	patient := models.Patient{
		Base: models.Base{
			ResourceType: models.ResourceTypePatient,
		},
	}
	if patient.ResourceType != models.ResourceTypePatient {
		t.Errorf("Expected ResourceType Patient, got %s", patient.ResourceType)
	}
}

func TestPatient_EmptyName(t *testing.T) {
	patient := models.Patient{
		Base: models.Base{
			ResourceType: models.ResourceTypePatient,
			ID:           "unit-2",
		},
		Active: ptrBool(true),
	}
	data, err := json.Marshal(patient)
	if err != nil {
		t.Fatalf("Failed to marshal patient with empty name: %v", err)
	}
	var out models.Patient
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("Failed to unmarshal patient with empty name: %v", err)
	}
	if out.ID != patient.ID {
		t.Errorf("ID mismatch: got %s, want %s", out.ID, patient.ID)
	}
	if out.Name != nil && len(out.Name) != 0 {
		t.Errorf("Expected empty Name, got %+v", out.Name)
	}
}

func TestPatient_InvalidBirthDate(t *testing.T) {
	jsonData := `{"resourceType":"Patient","id":"unit-3","birthDate":"not-a-date"}`
	var patient models.Patient
	err := json.Unmarshal([]byte(jsonData), &patient)
	if err == nil {
		t.Error("Expected error for invalid birthDate, got nil")
	}
}

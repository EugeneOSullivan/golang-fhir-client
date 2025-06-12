package tests

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/operations"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/search"
)

const hapiBaseURL = "https://hapi.fhir.org/baseR4"

func TestE2E_Patient_CRUD(t *testing.T) {
	if os.Getenv("SKIP_E2E") == "1" {
		t.Skip("Skipping E2E test due to SKIP_E2E env var")
	}

	ctx := context.Background()
	client := http.DefaultClient
	op := operations.NewHTTPOperation(client, hapiBaseURL)

	// 1. Create Patient
	patient := &models.Patient{
		Base: models.Base{
			ResourceType: models.ResourceTypePatient,
		},
		Active: ptrBool(true),
		Name: []models.HumanName{{
			Family: "E2ETest",
			Given:  []string{"Alice"},
		}},
		Gender:    "female",
		BirthDate: ptrTime(time.Date(1990, 2, 3, 0, 0, 0, 0, time.UTC)),
	}
	created, err := op.Create(ctx, string(models.ResourceTypePatient), patient)
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}
	createdPatient, ok := created.(*models.Patient)
	if !ok {
		t.Fatalf("Created resource is not a Patient")
	}
	if createdPatient.ID == "" {
		t.Fatalf("Created patient has no ID")
	}

	// 2. Read Patient
	read, err := op.Read(ctx, string(models.ResourceTypePatient), createdPatient.ID)
	if err != nil {
		t.Fatalf("Failed to read patient: %v", err)
	}
	readPatient, ok := read.(*models.Patient)
	if !ok {
		t.Fatalf("Read resource is not a Patient")
	}
	if readPatient.ID != createdPatient.ID {
		t.Errorf("Read patient ID mismatch: got %s, want %s", readPatient.ID, createdPatient.ID)
	}

	// 3. Search Patient
	params := search.NewParameters().Add("family", "E2ETest").Add("given", "Alice")
	bundle, err := op.Search(ctx, string(models.ResourceTypePatient), params)
	if err != nil {
		t.Fatalf("Failed to search patients: %v", err)
	}
	found := false
	for _, entry := range bundle.Entry {
		var p models.Patient
		if err := p.UnmarshalJSON(entry.Resource); err == nil && p.ID == createdPatient.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Created patient not found in search results")
	}

	// 4. Delete Patient
	if err := op.Delete(ctx, string(models.ResourceTypePatient), createdPatient.ID); err != nil {
		t.Errorf("Failed to delete patient: %v", err)
	}
}

func ptrBool(b bool) *bool           { return &b }
func ptrTime(t time.Time) *time.Time { return &t }

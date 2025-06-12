package examples_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/operations"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/search"
)

// setupTestClient creates a new FHIR client for testing
func setupTestClient(t *testing.T) *operations.HTTPOperation {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return operations.NewHTTPOperation(client, "http://hapi.fhir.org/baseR4")
}

// TestBasicCRUD demonstrates and tests basic CRUD operations on a Patient resource.
// It creates a new patient, reads it back, updates it, and finally deletes it.
func TestBasicCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	op := setupTestClient(t)
	ctx := context.Background()

	// Create a new patient
	newPatient := models.NewPatient()
	newPatient.Active = new(bool)
	*newPatient.Active = true
	newPatient.Name = []models.HumanName{
		{
			Family: "TestPatient",
			Given:  []string{"Integration", "Test"},
		},
	}
	newPatient.Gender = "male"
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	newPatient.BirthDate = &birthDate

	// Create
	createdResource, err := op.Create(ctx, "Patient", newPatient)
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}

	createdPatient, ok := createdResource.(*models.Patient)
	if !ok {
		t.Fatal("Created resource is not a Patient")
	}

	// Read
	readResource, err := op.Read(ctx, "Patient", createdPatient.ID)
	if err != nil {
		t.Fatalf("Failed to read patient: %v", err)
	}

	readPatient, ok := readResource.(*models.Patient)
	if !ok {
		t.Fatal("Read resource is not a Patient")
	}

	if readPatient.Name[0].Family != "TestPatient" {
		t.Errorf("Expected family name TestPatient, got %s", readPatient.Name[0].Family)
	}

	// Update
	readPatient.Name[0].Given = append(readPatient.Name[0].Given, "Updated")
	updatedResource, err := op.Update(ctx, "Patient", readPatient.ID, readPatient)
	if err != nil {
		t.Fatalf("Failed to update patient: %v", err)
	}

	updatedPatient, ok := updatedResource.(*models.Patient)
	if !ok {
		t.Fatal("Updated resource is not a Patient")
	}

	if len(updatedPatient.Name[0].Given) != 3 {
		t.Errorf("Expected 3 given names, got %d", len(updatedPatient.Name[0].Given))
	}

	// Delete
	err = op.Delete(ctx, "Patient", updatedPatient.ID)
	if err != nil {
		t.Fatalf("Failed to delete patient: %v", err)
	}

	// Verify deletion
	_, err = op.Read(ctx, "Patient", updatedPatient.ID)
	if err == nil {
		t.Error("Expected error when reading deleted patient")
	}
}

// TestSearch demonstrates and tests the search functionality.
// It searches for patients and verifies the search results.
func TestSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	op := setupTestClient(t)
	ctx := context.Background()

	// Create search parameters
	params := search.NewParameters()
	params.Add("name", "Test")
	params.Add("_count", "5")

	// Search
	bundle, err := op.Search(ctx, "Patient", params)
	if err != nil {
		t.Fatalf("Failed to search patients: %v", err)
	}

	if bundle.Total == nil {
		t.Fatal("Bundle total is nil")
	}

	// Process results
	for _, entry := range bundle.Entry {
		resource, err := bundle.GetTypedResource(entry.Resource)
		if err != nil {
			t.Errorf("Failed to convert resource: %v", err)
			continue
		}

		patient, ok := resource.(*models.Patient)
		if !ok {
			t.Error("Resource is not a Patient")
			continue
		}

		// Verify patient has required fields
		if patient.ID == "" {
			t.Error("Patient ID is empty")
		}
	}
}

// TestHistory demonstrates and tests the history functionality.
// It creates a patient, updates it to create multiple versions,
// and then retrieves and verifies the version history.
func TestHistory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	op := setupTestClient(t)
	ctx := context.Background()

	// Create a patient that we'll modify
	newPatient := models.NewPatient()
	newPatient.Active = new(bool)
	*newPatient.Active = true
	newPatient.Name = []models.HumanName{
		{
			Family: "HistoryTest",
			Given:  []string{"Version", "One"},
		},
	}

	// Create initial version
	createdResource, err := op.Create(ctx, "Patient", newPatient)
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}

	createdPatient, ok := createdResource.(*models.Patient)
	if !ok {
		t.Fatal("Created resource is not a Patient")
	}

	// Update to create a new version
	createdPatient.Name[0].Given = []string{"Version", "Two"}
	_, err = op.Update(ctx, "Patient", createdPatient.ID, createdPatient)
	if err != nil {
		t.Fatalf("Failed to update patient: %v", err)
	}

	// Get history
	params := search.NewParameters()
	params.Add("_count", "2")

	historyBundle, err := op.History(ctx, "Patient", createdPatient.ID, params)
	if err != nil {
		t.Fatalf("Failed to get history: %v", err)
	}

	if len(historyBundle.Entry) < 2 {
		t.Errorf("Expected at least 2 history entries, got %d", len(historyBundle.Entry))
	}

	// Clean up
	err = op.Delete(ctx, "Patient", createdPatient.ID)
	if err != nil {
		t.Fatalf("Failed to delete patient: %v", err)
	}
}

// TestVersioning demonstrates and tests version-specific operations.
// It creates a patient, updates it, and then retrieves specific versions
// using the vread operation.
func TestVersioning(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	op := setupTestClient(t)
	ctx := context.Background()

	// Create initial version
	newPatient := models.NewPatient()
	newPatient.Name = []models.HumanName{
		{
			Family: "VersionTest",
			Given:  []string{"Initial"},
		},
	}

	createdResource, err := op.Create(ctx, "Patient", newPatient)
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}

	createdPatient, ok := createdResource.(*models.Patient)
	if !ok {
		t.Fatal("Created resource is not a Patient")
	}

	// Get the version ID
	if createdPatient.Meta == nil || createdPatient.Meta.VersionID == "" {
		t.Fatal("Created patient has no version ID")
	}
	initialVersionId := createdPatient.Meta.VersionID

	// Update to create a new version
	createdPatient.Name[0].Given = []string{"Updated"}
	updatedResource, err := op.Update(ctx, "Patient", createdPatient.ID, createdPatient)
	if err != nil {
		t.Fatalf("Failed to update patient: %v", err)
	}

	updatedPatient, ok := updatedResource.(*models.Patient)
	if !ok {
		t.Fatal("Updated resource is not a Patient")
	}

	// Verify updated version has a different version ID
	if updatedPatient.Meta == nil || updatedPatient.Meta.VersionID == initialVersionId {
		t.Error("Updated patient should have a different version ID")
	}

	// Get specific version
	vreadResource, err := op.Vread(ctx, "Patient", createdPatient.ID, initialVersionId)
	if err != nil {
		t.Fatalf("Failed to read specific version: %v", err)
	}

	vreadPatient, ok := vreadResource.(*models.Patient)
	if !ok {
		t.Fatal("Vread resource is not a Patient")
	}

	// Verify it's the initial version
	if vreadPatient.Name[0].Given[0] != "Initial" {
		t.Errorf("Expected initial version name 'Initial', got %s", vreadPatient.Name[0].Given[0])
	}

	// Clean up
	err = op.Delete(ctx, "Patient", createdPatient.ID)
	if err != nil {
		t.Fatalf("Failed to delete patient: %v", err)
	}
}

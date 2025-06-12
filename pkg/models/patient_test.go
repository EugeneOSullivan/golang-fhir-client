package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPatientUnmarshalJSON(t *testing.T) {
	jsonData := `{
		"resourceType": "Patient",
		"id": "123",
		"active": true,
		"name": [
			{
				"family": "Doe",
				"given": ["John", "Middle"]
			}
		],
		"gender": "male",
		"birthDate": "2000-01-01",
		"deceasedDateTime": "2023-01-01T12:00:00Z",
		"address": [
			{
				"use": "home",
				"line": ["123 Main St"],
				"city": "Anytown",
				"state": "ST",
				"postalCode": "12345"
			}
		]
	}`

	var patient Patient
	err := json.Unmarshal([]byte(jsonData), &patient)
	if err != nil {
		t.Fatalf("Failed to unmarshal patient: %v", err)
	}

	// Verify resource type
	if patient.ResourceType != ResourceTypePatient {
		t.Errorf("Expected resource type Patient, got %s", patient.ResourceType)
	}

	// Verify ID
	if patient.ID != "123" {
		t.Errorf("Expected ID 123, got %s", patient.ID)
	}

	// Verify active status
	if patient.Active == nil || !*patient.Active {
		t.Error("Expected active to be true")
	}

	// Verify name
	if len(patient.Name) != 1 {
		t.Fatalf("Expected 1 name, got %d", len(patient.Name))
	}
	if patient.Name[0].Family != "Doe" {
		t.Errorf("Expected family name Doe, got %s", patient.Name[0].Family)
	}
	if len(patient.Name[0].Given) != 2 {
		t.Fatalf("Expected 2 given names, got %d", len(patient.Name[0].Given))
	}
	if patient.Name[0].Given[0] != "John" {
		t.Errorf("Expected first given name John, got %s", patient.Name[0].Given[0])
	}

	// Verify gender
	if patient.Gender != "male" {
		t.Errorf("Expected gender male, got %s", patient.Gender)
	}

	// Verify birth date
	expectedBirthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if patient.BirthDate == nil || !patient.BirthDate.Equal(expectedBirthDate) {
		t.Errorf("Expected birth date %v, got %v", expectedBirthDate, patient.BirthDate)
	}

	// Verify deceased date
	expectedDeceasedDate := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	if patient.DeceasedAt == nil || !patient.DeceasedAt.Equal(expectedDeceasedDate) {
		t.Errorf("Expected deceased date %v, got %v", expectedDeceasedDate, patient.DeceasedAt)
	}

	// Verify address
	if len(patient.Address) != 1 {
		t.Fatalf("Expected 1 address, got %d", len(patient.Address))
	}
	addr := patient.Address[0]
	if addr.Use != "home" {
		t.Errorf("Expected address use home, got %s", addr.Use)
	}
	if len(addr.Line) != 1 || addr.Line[0] != "123 Main St" {
		t.Errorf("Expected address line '123 Main St', got %v", addr.Line)
	}
	if addr.City != "Anytown" {
		t.Errorf("Expected city Anytown, got %s", addr.City)
	}
	if addr.State != "ST" {
		t.Errorf("Expected state ST, got %s", addr.State)
	}
	if addr.PostalCode != "12345" {
		t.Errorf("Expected postal code 12345, got %s", addr.PostalCode)
	}
}

func TestBundleUnmarshalWithTypedResources(t *testing.T) {
	jsonData := `{
		"resourceType": "Bundle",
		"type": "searchset",
		"total": 1,
		"entry": [
			{
				"fullUrl": "http://example.com/Patient/123",
				"resource": {
					"resourceType": "Patient",
					"id": "123",
					"name": [
						{
							"family": "Doe",
							"given": ["John"]
						}
					],
					"birthDate": "2000-01-01"
				}
			}
		]
	}`

	var bundle Bundle
	err := json.Unmarshal([]byte(jsonData), &bundle)
	if err != nil {
		t.Fatalf("Failed to unmarshal bundle: %v", err)
	}

	// Verify bundle metadata
	if bundle.Type != "searchset" {
		t.Errorf("Expected bundle type searchset, got %s", bundle.Type)
	}
	if bundle.Total == nil || *bundle.Total != 1 {
		t.Error("Expected total to be 1")
	}

	// Verify entries
	if len(bundle.Entry) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(bundle.Entry))
	}

	entry := bundle.Entry[0]
	if entry.FullURL != "http://example.com/Patient/123" {
		t.Errorf("Expected fullUrl http://example.com/Patient/123, got %s", entry.FullURL)
	}

	// Convert raw resource to typed resource
	resource, err := bundle.GetTypedResource(entry.Resource)
	if err != nil {
		t.Fatalf("Failed to get typed resource: %v", err)
	}

	// Type assert to Patient
	patient, ok := resource.(*Patient)
	if !ok {
		t.Fatal("Expected resource to be a Patient")
	}

	// Verify patient data
	if patient.ID != "123" {
		t.Errorf("Expected patient ID 123, got %s", patient.ID)
	}
	if len(patient.Name) != 1 {
		t.Fatalf("Expected 1 name, got %d", len(patient.Name))
	}
	if patient.Name[0].Family != "Doe" {
		t.Errorf("Expected family name Doe, got %s", patient.Name[0].Family)
	}
	if len(patient.Name[0].Given) != 1 || patient.Name[0].Given[0] != "John" {
		t.Errorf("Expected given name John, got %v", patient.Name[0].Given)
	}

	expectedBirthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if patient.BirthDate == nil || !patient.BirthDate.Equal(expectedBirthDate) {
		t.Errorf("Expected birth date %v, got %v", expectedBirthDate, patient.BirthDate)
	}
}

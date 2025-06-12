package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/operations"
)

func main() {
	// Create a new HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create a new FHIR client
	op := operations.NewHTTPOperation(client, "http://hapi.fhir.org/baseR4")

	// Create a new patient
	newPatient := models.NewPatient()
	newPatient.Active = new(bool)
	*newPatient.Active = true
	newPatient.Name = []models.HumanName{
		{
			Family: "Doe",
			Given:  []string{"John"},
		},
	}
	newPatient.Gender = "male"
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	newPatient.BirthDate = &birthDate

	// Create the patient
	ctx := context.Background()
	createdResource, err := op.Create(ctx, "Patient", newPatient)
	if err != nil {
		log.Fatalf("Failed to create patient: %v", err)
	}

	// Type assert to Patient
	createdPatient, ok := createdResource.(*models.Patient)
	if !ok {
		log.Fatal("Created resource is not a Patient")
	}

	fmt.Printf("Created patient with ID: %s\n", createdPatient.ID)

	// Read the patient back
	readResource, err := op.Read(ctx, "Patient", createdPatient.ID)
	if err != nil {
		log.Fatalf("Failed to read patient: %v", err)
	}

	// Type assert to Patient
	readPatient, ok := readResource.(*models.Patient)
	if !ok {
		log.Fatal("Read resource is not a Patient")
	}

	fmt.Printf("Read patient name: %s, %s\n",
		readPatient.Name[0].Family,
		readPatient.Name[0].Given[0])

	// Update the patient
	readPatient.Name[0].Given = append(readPatient.Name[0].Given, "Middle")
	updatedResource, err := op.Update(ctx, "Patient", readPatient.ID, readPatient)
	if err != nil {
		log.Fatalf("Failed to update patient: %v", err)
	}

	// Type assert to Patient
	updatedPatient, ok := updatedResource.(*models.Patient)
	if !ok {
		log.Fatal("Updated resource is not a Patient")
	}

	fmt.Printf("Updated patient with middle name: %v\n", updatedPatient.Name[0].Given)

	// Delete the patient
	err = op.Delete(ctx, "Patient", updatedPatient.ID)
	if err != nil {
		log.Fatalf("Failed to delete patient: %v", err)
	}

	fmt.Printf("Deleted patient with ID: %s\n", updatedPatient.ID)
}

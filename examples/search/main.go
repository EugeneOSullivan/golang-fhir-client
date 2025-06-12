package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/operations"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/search"
)

func main() {
	// Create a new HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create a new FHIR client
	op := operations.NewHTTPOperation(client, "http://hapi.fhir.org/baseR4")

	// Create search parameters
	params := search.NewParameters()
	params.Add("name", "Smith")       // Search for patients with name containing "Smith"
	params.Add("birthdate", "ge2000") // Born after year 2000
	params.Add("_count", "5")         // Limit to 5 results

	// Search for patients
	ctx := context.Background()
	bundle, err := op.Search(ctx, "Patient", params)
	if err != nil {
		log.Fatalf("Failed to search patients: %v", err)
	}

	// Print total results
	if bundle.Total != nil {
		fmt.Printf("Found %d patients\n", *bundle.Total)
	}

	// Process each entry in the bundle
	for i, entry := range bundle.Entry {
		// Convert the raw resource to a Patient
		resource, err := bundle.GetTypedResource(entry.Resource)
		if err != nil {
			log.Printf("Failed to convert resource %d: %v", i, err)
			continue
		}

		// Type assert to Patient
		patient, ok := resource.(*models.Patient)
		if !ok {
			log.Printf("Resource %d is not a Patient", i)
			continue
		}

		// Print patient information
		fmt.Printf("\nPatient %d:\n", i+1)
		fmt.Printf("ID: %s\n", patient.ID)

		// Print names
		for _, name := range patient.Name {
			if name.Family != "" {
				fmt.Printf("Name: %s", name.Family)
				if len(name.Given) > 0 {
					fmt.Printf(", %s", name.Given[0])
				}
				fmt.Println()
			}
		}

		// Print birth date
		if patient.BirthDate != nil {
			fmt.Printf("Birth Date: %s\n", patient.BirthDate.Format("2006-01-02"))
		}

		// Print gender
		if patient.Gender != "" {
			fmt.Printf("Gender: %s\n", patient.Gender)
		}

		// Print contact points
		for _, telecom := range patient.Telecom {
			if telecom.System != "" && telecom.Value != "" {
				fmt.Printf("Contact (%s): %s\n", telecom.System, telecom.Value)
			}
		}

		// Print addresses
		for _, address := range patient.Address {
			if len(address.Line) > 0 || address.City != "" {
				fmt.Println("Address:")
				for _, line := range address.Line {
					fmt.Printf("  %s\n", line)
				}
				if address.City != "" {
					fmt.Printf("  %s", address.City)
					if address.State != "" {
						fmt.Printf(", %s", address.State)
					}
					if address.PostalCode != "" {
						fmt.Printf(" %s", address.PostalCode)
					}
					fmt.Println()
				}
			}
		}
	}
}

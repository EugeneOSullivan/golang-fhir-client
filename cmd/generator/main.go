package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type StructureDefinition struct {
	ResourceType string `json:"resourceType"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Snapshot     struct {
		Element []Element `json:"element"`
	} `json:"snapshot"`
}

type Element struct {
	Path       string `json:"path"`
	Min        int    `json:"min"`
	Max        string `json:"max"`
	Type       []Type `json:"type"`
	Definition string `json:"definition"`
}

type Type struct {
	Code string `json:"code"`
}

func main() {
	var (
		inputDir  = flag.String("input", "", "Directory containing FHIR StructureDefinitions")
		outputDir = flag.String("output", "pkg/models", "Output directory for generated Go files")
	)
	flag.Parse()

	if *inputDir == "" {
		log.Fatal("Input directory is required")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// TODO: Implement the actual generation logic
	fmt.Println("FHIR model generator initialized")
	fmt.Printf("Input directory: %s\n", *inputDir)
	fmt.Printf("Output directory: %s\n", *outputDir)
}

// Helper function to convert FHIR types to Go types
func fhirTypeToGoType(fhirType string) string {
	switch fhirType {
	case "boolean":
		return "bool"
	case "integer":
		return "int"
	case "decimal":
		return "float64"
	case "string", "code", "id", "uri", "url", "canonical", "markdown":
		return "string"
	case "date", "dateTime", "instant", "time":
		return "time.Time"
	default:
		return "interface{}"
	}
}

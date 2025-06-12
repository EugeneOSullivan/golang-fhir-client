# Golang FHIR Client

A robust, type-safe FHIR client implementation in Go, supporting both R4 and R5 versions of the FHIR specification.

## Features

- Type-safe FHIR resource handling
- Support for FHIR R4 and R5
- Automatic JSON marshaling/unmarshaling
- Comprehensive CRUD operations
- Advanced search capabilities
- Bundle support
- Version-aware client operations
- Cross-version resource mapping

## Installation

```bash
go get github.com/eugeneosullivan/golang-fhir-client
```

## Quick Start

```go
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
    // Create a new client
    client := operations.NewHTTPOperation(
        &http.Client{Timeout: 30 * time.Second},
        "http://hapi.fhir.org/baseR4",
    )

    // Search for patients
    params := search.NewParameters()
    params.Add("name", "Smith")
    
    ctx := context.Background()
    bundle, err := client.Search(ctx, "Patient", params)
    if err != nil {
        log.Fatal(err)
    }

    // Process results
    for _, entry := range bundle.Entry {
        resource, err := bundle.GetTypedResource(entry.Resource)
        if err != nil {
            log.Printf("Error: %v", err)
            continue
        }

        patient, ok := resource.(*models.Patient)
        if !ok {
            continue
        }

        fmt.Printf("Found patient: %s\n", patient.ID)
    }
}
```

## Project Structure

```
golang-fhir-client/
├── pkg/
│   ├── models/         # Base resource models
│   │   ├── r4/        # R4-specific resource definitions
│   │   └── r5/        # R5-specific resource definitions
│   ├── operations/     # FHIR operations implementation
│   └── search/         # Search parameter handling
├── examples/           # Usage examples
└── tests/             # Integration tests
```

## Version Support

The client supports both FHIR R4 and R5 versions through version-specific packages:

- `models/r4`: Contains R4-specific resource definitions and validations
- `models/r5`: Contains R5-specific resource definitions and validations

The base models in `models/` package contain common fields and functionality shared between versions.

## Supported Operations

- Read: Get a specific resource by ID
- VRead: Get a specific version of a resource
- Create: Create a new resource
- Update: Update an existing resource
- Delete: Delete a resource
- Search: Search for resources with parameters
- History: Get resource version history
- Transaction: Execute a batch of operations
- Operation: Execute custom operations

## Search Parameters

The search package provides a fluent interface for building search queries:

```go
params := search.NewParameters()
params.Add("name", "Smith")
params.Add("birthdate", "ge2000")
params.Add("_count", "5")
```

## Error Handling

The client provides detailed error information for:

- HTTP errors
- JSON marshaling/unmarshaling errors
- FHIR operation outcomes
- Validation errors

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [HL7 FHIR](https://www.hl7.org/fhir/) - FHIR specification
- [HAPI FHIR](https://hapifhir.io/) - Reference server used for testing 

## Testing

This project includes both unit and end-to-end (E2E) tests for FHIR resources, especially the Patient resource.

- **Unit tests** are located in the `tests/` directory (e.g., `patient_unit_test.go`). These cover JSON marshaling/unmarshaling, required fields, and edge cases for Patient.
- **E2E tests** (e.g., `e2e_patient_test.go`) interact with a live HAPI FHIR server (https://hapi.fhir.org/baseR4) to test create, read, search, and delete operations for Patient resources.

### Running Tests

To run all tests in the `tests` directory:

```bash
go test ./tests/... -v
```

To skip E2E tests (for example, in CI), set the environment variable:

```bash
SKIP_E2E=1 go test ./tests/... -v
```

You can also run all tests (including any in the main codebase) with:

```bash
go test ./... ./tests/... -v
``` 
package operations

import (
	"context"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/search"
)

// Operation defines the interface for FHIR operations
type Operation interface {
	// Read retrieves a resource by ID
	Read(ctx context.Context, resourceType, id string) (models.Resource, error)

	// Vread retrieves a specific version of a resource
	Vread(ctx context.Context, resourceType, id, versionId string) (models.Resource, error)

	// Create creates a new resource
	Create(ctx context.Context, resourceType string, resource interface{}) (models.Resource, error)

	// Update updates an existing resource
	Update(ctx context.Context, resourceType, id string, resource interface{}) (models.Resource, error)

	// Patch patches an existing resource
	Patch(ctx context.Context, resourceType, id string, patchBody interface{}) (models.Resource, error)

	// Delete deletes a resource
	Delete(ctx context.Context, resourceType, id string) error

	// Search searches for resources
	Search(ctx context.Context, resourceType string, params *search.Parameters) (*models.Bundle, error)

	// History gets the history of a resource
	History(ctx context.Context, resourceType, id string, params *search.Parameters) (*models.Bundle, error)

	// Transaction executes a batch of operations
	Transaction(ctx context.Context, bundle interface{}) (*models.Bundle, error)

	// Capabilities retrieves the server's capability statement
	Capabilities(ctx context.Context) (models.Resource, error)

	// Operation executes a custom operation
	Operation(ctx context.Context, name string, input interface{}) (models.Resource, error)
}

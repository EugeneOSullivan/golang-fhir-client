package operations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/models"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/search"
)

// HTTPOperation implements the Operation interface using HTTP
type HTTPOperation struct {
	client  *http.Client
	baseURL string
	headers map[string]string
	mapper  *models.ResourceMapper
}

// NewHTTPOperation creates a new HTTP operation handler
func NewHTTPOperation(client *http.Client, baseURL string) *HTTPOperation {
	if client == nil {
		client = http.DefaultClient
	}

	return &HTTPOperation{
		client:  client,
		baseURL: baseURL,
		headers: make(map[string]string),
		mapper:  models.NewResourceMapper(),
	}
}

// doRequest performs an HTTP request and returns the response
func (o *HTTPOperation) doRequest(ctx context.Context, method, url string, body interface{}) (json.RawMessage, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/fhir+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/fhir+json")
	}
	for k, v := range o.headers {
		req.Header.Set(k, v)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned error status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Read retrieves a resource by ID and returns a typed resource
func (o *HTTPOperation) Read(ctx context.Context, resourceType, id string) (models.Resource, error) {
	url := o.buildURL(resourceType, id)
	data, err := o.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalResource(data)
}

// Vread retrieves a specific version of a resource and returns a typed resource
func (o *HTTPOperation) Vread(ctx context.Context, resourceType, id, versionId string) (models.Resource, error) {
	url := o.buildURL(resourceType, id, "_history", versionId)
	data, err := o.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalResource(data)
}

// Create creates a new resource and returns the typed created resource
func (o *HTTPOperation) Create(ctx context.Context, resourceType string, resource interface{}) (models.Resource, error) {
	url := o.buildURL(resourceType)
	data, err := o.doRequest(ctx, http.MethodPost, url, resource)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalResource(data)
}

// Update updates an existing resource and returns the typed updated resource
func (o *HTTPOperation) Update(ctx context.Context, resourceType, id string, resource interface{}) (models.Resource, error) {
	url := o.buildURL(resourceType, id)
	data, err := o.doRequest(ctx, http.MethodPut, url, resource)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalResource(data)
}

// Patch patches an existing resource and returns the typed patched resource
func (o *HTTPOperation) Patch(ctx context.Context, resourceType, id string, patchBody interface{}) (models.Resource, error) {
	url := o.buildURL(resourceType, id)
	data, err := o.doRequest(ctx, http.MethodPatch, url, patchBody)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalResource(data)
}

// Delete deletes a resource
func (o *HTTPOperation) Delete(ctx context.Context, resourceType, id string) error {
	url := o.buildURL(resourceType, id)
	_, err := o.doRequest(ctx, http.MethodDelete, url, nil)
	return err
}

// Search searches for resources and returns a typed Bundle
func (o *HTTPOperation) Search(ctx context.Context, resourceType string, params *search.Parameters) (*models.Bundle, error) {
	url := o.buildURL(resourceType)
	if params != nil {
		url += "?" + params.Encode()
	}
	data, err := o.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalBundle(data)
}

// History gets the history of a resource and returns a typed Bundle
func (o *HTTPOperation) History(ctx context.Context, resourceType, id string, params *search.Parameters) (*models.Bundle, error) {
	url := o.buildURL(resourceType, id, "_history")
	if params != nil {
		url += "?" + params.Encode()
	}
	data, err := o.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalBundle(data)
}

// Transaction executes a batch of operations and returns a typed Bundle
func (o *HTTPOperation) Transaction(ctx context.Context, bundle interface{}) (*models.Bundle, error) {
	url := o.buildURL()
	data, err := o.doRequest(ctx, http.MethodPost, url, bundle)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalBundle(data)
}

// Capabilities retrieves the server's capability statement as a typed Resource
func (o *HTTPOperation) Capabilities(ctx context.Context) (models.Resource, error) {
	url := o.buildURL("metadata")
	data, err := o.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalResource(data)
}

// Operation executes a custom operation and returns a typed Resource
func (o *HTTPOperation) Operation(ctx context.Context, name string, input interface{}) (models.Resource, error) {
	url := o.buildURL("$" + name)
	data, err := o.doRequest(ctx, http.MethodPost, url, input)
	if err != nil {
		return nil, err
	}
	return o.mapper.UnmarshalResource(data)
}

// buildURL constructs the full URL for a FHIR request
func (o *HTTPOperation) buildURL(parts ...string) string {
	url := o.baseURL
	for _, part := range parts {
		if part != "" {
			url += "/" + part
		}
	}
	return url
}

// SetHeader sets a custom header for all operations
func (o *HTTPOperation) SetHeader(key, value string) {
	o.headers[key] = value
}

package client

import (
	"context"
	"net/http"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/mapper"
	"github.com/eugeneosullivan/golang-fhir-client/pkg/version"
)

// Config holds the FHIR client configuration
type Config struct {
	BaseURL     string
	HTTPClient  *http.Client
	AuthConfig  *AuthConfig
	FHIRVersion version.Version
}

// AuthConfig holds OAuth2 configuration
type AuthConfig struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
}

// Client represents a FHIR client
type Client struct {
	config         *Config
	httpClient     *http.Client
	versionManager version.VersionManager
	mapper         mapper.Mapper
}

// NewClient creates a new FHIR client with the given configuration
func NewClient(config *Config) (*Client, error) {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	if config.FHIRVersion == "" {
		config.FHIRVersion = version.R4 // Default to R4
	}

	versionManager, err := version.NewVersionManager(config.FHIRVersion)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:         config,
		httpClient:     config.HTTPClient,
		versionManager: versionManager,
		mapper:         mapper.NewMapper(),
	}, nil
}

// Config returns the client configuration
func (c *Client) Config() *Config {
	return c.config
}

// Version returns the FHIR version manager
func (c *Client) Version() version.VersionManager {
	return c.versionManager
}

// MapResource converts a resource from the client's version to the target version
func (c *Client) MapResource(resource interface{}, targetVersion version.Version) (interface{}, error) {
	return c.mapper.MapResource(resource, c.config.FHIRVersion, targetVersion)
}

// Resource represents a generic FHIR resource interface
type Resource interface {
	GetResourceType() string
	GetVersion() version.Version
}

// Operation represents a FHIR operation
type Operation struct {
	client  *Client
	ctx     context.Context
	version version.Version // Target version for this operation
}

// WithContext sets the context for the operation
func (o *Operation) WithContext(ctx context.Context) *Operation {
	o.ctx = ctx
	return o
}

// WithVersion sets the target version for this operation
func (o *Operation) WithVersion(v version.Version) *Operation {
	o.version = v
	return o
}

// NewOperation creates a new operation
func (c *Client) NewOperation() *Operation {
	return &Operation{
		client:  c,
		ctx:     context.Background(),
		version: c.config.FHIRVersion,
	}
}

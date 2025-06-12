package version

import "fmt"

// Version represents a FHIR version
type Version string

const (
	R4 Version = "R4"
	R5 Version = "R5"
)

// VersionInfo contains metadata about a FHIR version
type VersionInfo struct {
	Version     Version
	APIVersion  string
	BaseURL     string
	Conformance string
}

// VersionManager handles version-specific operations
type VersionManager interface {
	// GetVersion returns the FHIR version
	GetVersion() Version

	// GetVersionInfo returns version-specific information
	GetVersionInfo() *VersionInfo

	// IsSupported checks if a specific resource type is supported in this version
	IsSupported(resourceType string) bool

	// GetBaseResource returns the base resource type for version-specific resources
	GetBaseResource(resourceType string) string
}

// baseVersionManager provides common functionality for version managers
type baseVersionManager struct {
	version Version
	info    *VersionInfo
}

// R4Manager implements VersionManager for FHIR R4
type R4Manager struct {
	baseVersionManager
}

// R5Manager implements VersionManager for FHIR R5
type R5Manager struct {
	baseVersionManager
}

func newR4Manager() *R4Manager {
	return &R4Manager{
		baseVersionManager: baseVersionManager{
			version: R4,
			info: &VersionInfo{
				Version:     R4,
				APIVersion:  "4.0.1",
				BaseURL:     "/baseR4",
				Conformance: "metadata",
			},
		},
	}
}

func newR5Manager() *R5Manager {
	return &R5Manager{
		baseVersionManager: baseVersionManager{
			version: R5,
			info: &VersionInfo{
				Version:     R5,
				APIVersion:  "5.0.0",
				BaseURL:     "/R5",
				Conformance: "metadata",
			},
		},
	}
}

// GetVersion returns the FHIR version
func (b *baseVersionManager) GetVersion() Version {
	return b.version
}

// GetVersionInfo returns version-specific information
func (b *baseVersionManager) GetVersionInfo() *VersionInfo {
	return b.info
}

// IsSupported checks if a specific resource type is supported in this version
func (r *R4Manager) IsSupported(resourceType string) bool {
	// TODO: Implement R4-specific resource type validation
	return true
}

// GetBaseResource returns the base resource type for R4-specific resources
func (r *R4Manager) GetBaseResource(resourceType string) string {
	// TODO: Implement R4-specific resource type mapping
	return resourceType
}

// IsSupported checks if a specific resource type is supported in this version
func (r *R5Manager) IsSupported(resourceType string) bool {
	// TODO: Implement R5-specific resource type validation
	return true
}

// GetBaseResource returns the base resource type for R5-specific resources
func (r *R5Manager) GetBaseResource(resourceType string) string {
	// TODO: Implement R5-specific resource type mapping
	return resourceType
}

// NewVersionManager creates a new version manager for the specified FHIR version
func NewVersionManager(version Version) (VersionManager, error) {
	switch version {
	case R4:
		return newR4Manager(), nil
	case R5:
		return newR5Manager(), nil
	default:
		return nil, fmt.Errorf("unsupported FHIR version: %s", version)
	}
}

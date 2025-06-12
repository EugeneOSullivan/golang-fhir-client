package mapper

import (
	"fmt"

	"github.com/eugeneosullivan/golang-fhir-client/pkg/version"
)

// Mapper provides functionality to map resources between different FHIR versions
type Mapper interface {
	// MapResource converts a resource from one version to another
	MapResource(resource interface{}, fromVersion, toVersion version.Version) (interface{}, error)

	// CanMap checks if mapping is possible between two versions for a given resource type
	CanMap(resourceType string, fromVersion, toVersion version.Version) bool

	// GetMappingPath returns the sequence of transformations needed to map between versions
	GetMappingPath(fromVersion, toVersion version.Version) ([]version.Version, error)
}

// DefaultMapper provides a default implementation of the Mapper interface
type DefaultMapper struct {
	versionManagers map[version.Version]version.VersionManager
}

// NewMapper creates a new mapper instance
func NewMapper() *DefaultMapper {
	mapper := &DefaultMapper{
		versionManagers: make(map[version.Version]version.VersionManager),
	}

	// Initialize version managers for supported versions
	r4Manager, _ := version.NewVersionManager(version.R4)
	r5Manager, _ := version.NewVersionManager(version.R5)

	mapper.versionManagers[version.R4] = r4Manager
	mapper.versionManagers[version.R5] = r5Manager

	return mapper
}

// MapResource converts a resource from one version to another
func (m *DefaultMapper) MapResource(resource interface{}, fromVersion, toVersion version.Version) (interface{}, error) {
	if fromVersion == toVersion {
		return resource, nil
	}

	path, err := m.GetMappingPath(fromVersion, toVersion)
	if err != nil {
		return nil, err
	}

	result := resource
	for i := 0; i < len(path)-1; i++ {
		currentVersion := path[i]
		nextVersion := path[i+1]

		// TODO: Implement actual mapping logic between versions
		result, err = m.mapBetweenVersions(result, currentVersion, nextVersion)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// CanMap checks if mapping is possible between two versions for a given resource type
func (m *DefaultMapper) CanMap(resourceType string, fromVersion, toVersion version.Version) bool {
	fromManager, fromExists := m.versionManagers[fromVersion]
	toManager, toExists := m.versionManagers[toVersion]

	if !fromExists || !toExists {
		return false
	}

	return fromManager.IsSupported(resourceType) && toManager.IsSupported(resourceType)
}

// GetMappingPath returns the sequence of transformations needed to map between versions
func (m *DefaultMapper) GetMappingPath(fromVersion, toVersion version.Version) ([]version.Version, error) {
	if fromVersion == toVersion {
		return []version.Version{fromVersion}, nil
	}

	// For now, we only support direct mapping between adjacent versions
	// This can be extended to support more complex mapping paths
	if fromVersion == version.R4 && toVersion == version.R5 {
		return []version.Version{version.R4, version.R5}, nil
	}
	if fromVersion == version.R5 && toVersion == version.R4 {
		return []version.Version{version.R5, version.R4}, nil
	}

	return nil, fmt.Errorf("no mapping path available from %s to %s", fromVersion, toVersion)
}

// mapBetweenVersions performs the actual mapping between two adjacent versions
func (m *DefaultMapper) mapBetweenVersions(resource interface{}, fromVersion, toVersion version.Version) (interface{}, error) {
	// TODO: Implement version-specific mapping logic
	// This would include:
	// 1. Type conversions
	// 2. Field mappings
	// 3. Handling deprecated fields
	// 4. Adding new required fields
	return resource, nil
}

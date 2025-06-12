package models

import (
	"encoding/json"
	"fmt"
)

// ResourceMapper handles conversion between JSON and FHIR resource structs
type ResourceMapper struct {
	typeRegistry map[ResourceType]func() Resource
}

// NewResourceMapper creates a new resource mapper with default resource types
func NewResourceMapper() *ResourceMapper {
	m := &ResourceMapper{
		typeRegistry: make(map[ResourceType]func() Resource),
	}

	// Register default resource types
	m.RegisterResource(ResourceTypePatient, func() Resource { return NewPatient() })
	// Add more resource types as they are implemented

	return m
}

// RegisterResource registers a new resource type with the mapper
func (m *ResourceMapper) RegisterResource(resourceType ResourceType, factory func() Resource) {
	m.typeRegistry[resourceType] = factory
}

// UnmarshalResource converts JSON data to a FHIR resource
func (m *ResourceMapper) UnmarshalResource(data []byte) (Resource, error) {
	// First unmarshal just the resource type
	var typeHolder struct {
		ResourceType ResourceType `json:"resourceType"`
	}

	if err := json.Unmarshal(data, &typeHolder); err != nil {
		return nil, fmt.Errorf("failed to determine resource type: %w", err)
	}

	// Get the factory function for this resource type
	factory, ok := m.typeRegistry[typeHolder.ResourceType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", typeHolder.ResourceType)
	}

	// Create a new instance of the resource
	resource := factory()

	// Unmarshal the full data into the resource
	if err := json.Unmarshal(data, resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	return resource, nil
}

// UnmarshalBundle converts a FHIR Bundle JSON to a Bundle struct with typed resources
func (m *ResourceMapper) UnmarshalBundle(data []byte) (*Bundle, error) {
	var bundle Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bundle: %w", err)
	}

	// Process each entry in the bundle
	for i, entry := range bundle.Entry {
		if entry.Resource == nil {
			continue
		}

		// Convert the raw resource to a typed resource
		resource, err := m.UnmarshalResource(entry.Resource)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal bundle entry %d: %w", i, err)
		}

		// Re-marshal the typed resource
		typedJSON, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal typed resource %d: %w", i, err)
		}

		bundle.Entry[i].Resource = typedJSON
	}

	return &bundle, nil
}

// MarshalResource converts a FHIR resource to JSON
func (m *ResourceMapper) MarshalResource(resource Resource) ([]byte, error) {
	return json.Marshal(resource)
}

// GetTypedResource attempts to convert a raw JSON resource to its typed struct
func (m *ResourceMapper) GetTypedResource(data json.RawMessage) (Resource, error) {
	if data == nil {
		return nil, nil
	}
	return m.UnmarshalResource(data)
}

// GetTypedResources converts a slice of raw JSON resources to typed structs
func (m *ResourceMapper) GetTypedResources(data []json.RawMessage) ([]Resource, error) {
	resources := make([]Resource, 0, len(data))

	for i, rawResource := range data {
		resource, err := m.GetTypedResource(rawResource)
		if err != nil {
			return nil, fmt.Errorf("failed to convert resource %d: %w", i, err)
		}
		if resource != nil {
			resources = append(resources, resource)
		}
	}

	return resources, nil
}

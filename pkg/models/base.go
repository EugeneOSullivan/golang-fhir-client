package models

import (
	"encoding/json"
	"time"
)

// Resource represents the base interface that all FHIR resources must implement
type Resource interface {
	GetResourceType() string
}

// ResourceType represents a FHIR resource type
type ResourceType string

// Base FHIR resource types
const (
	ResourceTypePatient           ResourceType = "Patient"
	ResourceTypeObservation       ResourceType = "Observation"
	ResourceTypeCondition         ResourceType = "Condition"
	ResourceTypeMedicationRequest ResourceType = "MedicationRequest"
	// Add more resource types as needed
)

// Base represents common fields present in all FHIR resources
type Base struct {
	ResourceType ResourceType      `json:"resourceType"`
	ID           string            `json:"id,omitempty"`
	Meta         *Meta             `json:"meta,omitempty"`
	Language     string            `json:"language,omitempty"`
	Text         *Narrative        `json:"text,omitempty"`
	Extension    []Extension       `json:"extension,omitempty"`
	Contained    []json.RawMessage `json:"contained,omitempty"`
}

// GetResourceType implements the Resource interface
func (b Base) GetResourceType() string {
	return string(b.ResourceType)
}

// Meta represents FHIR resource metadata
type Meta struct {
	VersionID   string     `json:"versionId,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Source      string     `json:"source,omitempty"`
	Profile     []string   `json:"profile,omitempty"`
	Security    []Coding   `json:"security,omitempty"`
	Tag         []Coding   `json:"tag,omitempty"`
}

// Narrative represents the human-readable summary of a resource
type Narrative struct {
	Status NarrativeStatus `json:"status"`
	Div    string          `json:"div"`
}

// NarrativeStatus represents the status of a narrative text
type NarrativeStatus string

const (
	NarrativeStatusGenerated  NarrativeStatus = "generated"
	NarrativeStatusExtensions NarrativeStatus = "extensions"
	NarrativeStatusAdditional NarrativeStatus = "additional"
	NarrativeStatusEmpty      NarrativeStatus = "empty"
)

// Extension represents a FHIR extension
type Extension struct {
	URL                  string           `json:"url"`
	ValueString          *string          `json:"valueString,omitempty"`
	ValueInteger         *int             `json:"valueInteger,omitempty"`
	ValueBoolean         *bool            `json:"valueBoolean,omitempty"`
	ValueCode            *string          `json:"valueCode,omitempty"`
	ValueDateTime        *time.Time       `json:"valueDateTime,omitempty"`
	ValueQuantity        *Quantity        `json:"valueQuantity,omitempty"`
	ValueReference       *Reference       `json:"valueReference,omitempty"`
	ValueCodeableConcept *CodeableConcept `json:"valueCodeableConcept,omitempty"`
}

// Coding represents a code from a code system
type Coding struct {
	System       string `json:"system,omitempty"`
	Version      string `json:"version,omitempty"`
	Code         string `json:"code,omitempty"`
	Display      string `json:"display,omitempty"`
	UserSelected *bool  `json:"userSelected,omitempty"`
}

// CodeableConcept represents a concept that may be defined by one or more code systems
type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}

// Quantity represents a measured or measurable amount
type Quantity struct {
	Value      *float64 `json:"value,omitempty"`
	Comparator string   `json:"comparator,omitempty"`
	Unit       string   `json:"unit,omitempty"`
	System     string   `json:"system,omitempty"`
	Code       string   `json:"code,omitempty"`
}

// Reference represents a reference to another resource
type Reference struct {
	Reference string `json:"reference,omitempty"`
	Type      string `json:"type,omitempty"`
	Display   string `json:"display,omitempty"`
}

// Period represents a time period
type Period struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// HumanName represents a human name
type HumanName struct {
	Use    string   `json:"use,omitempty"`
	Text   string   `json:"text,omitempty"`
	Family string   `json:"family,omitempty"`
	Given  []string `json:"given,omitempty"`
	Prefix []string `json:"prefix,omitempty"`
	Suffix []string `json:"suffix,omitempty"`
	Period *Period  `json:"period,omitempty"`
}

// ContactPoint represents contact information
type ContactPoint struct {
	System ContactPointSystem `json:"system,omitempty"`
	Value  string             `json:"value,omitempty"`
	Use    ContactPointUse    `json:"use,omitempty"`
	Rank   *int               `json:"rank,omitempty"`
	Period *Period            `json:"period,omitempty"`
}

// ContactPointSystem represents the type of contact point
type ContactPointSystem string

const (
	ContactPointSystemPhone ContactPointSystem = "phone"
	ContactPointSystemFax   ContactPointSystem = "fax"
	ContactPointSystemEmail ContactPointSystem = "email"
	ContactPointSystemPager ContactPointSystem = "pager"
	ContactPointSystemURL   ContactPointSystem = "url"
	ContactPointSystemSMS   ContactPointSystem = "sms"
	ContactPointSystemOther ContactPointSystem = "other"
)

// ContactPointUse represents the purpose of the contact point
type ContactPointUse string

const (
	ContactPointUseHome   ContactPointUse = "home"
	ContactPointUseWork   ContactPointUse = "work"
	ContactPointUseTemp   ContactPointUse = "temp"
	ContactPointUseOld    ContactPointUse = "old"
	ContactPointUseMobile ContactPointUse = "mobile"
)

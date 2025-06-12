package models

import (
	"encoding/json"
	"time"
)

// Patient represents a FHIR Patient resource
type Patient struct {
	Base
	Active               *bool            `json:"active,omitempty"`
	Name                 []HumanName      `json:"name,omitempty"`
	Telecom              []ContactPoint   `json:"telecom,omitempty"`
	Gender               string           `json:"gender,omitempty"`
	BirthDate            *time.Time       `json:"birthDate,omitempty"`
	Deceased             *bool            `json:"deceasedBoolean,omitempty"`
	DeceasedAt           *time.Time       `json:"deceasedDateTime,omitempty"`
	Address              []Address        `json:"address,omitempty"`
	MaritalStatus        *CodeableConcept `json:"maritalStatus,omitempty"`
	MultipleBirth        *bool            `json:"multipleBirthBoolean,omitempty"`
	MultipleBirthInt     *int             `json:"multipleBirthInteger,omitempty"`
	Photo                []Attachment     `json:"photo,omitempty"`
	Contact              []PatientContact `json:"contact,omitempty"`
	Communication        []Communication  `json:"communication,omitempty"`
	GeneralPractitioner  []Reference      `json:"generalPractitioner,omitempty"`
	ManagingOrganization *Reference       `json:"managingOrganization,omitempty"`
	Link                 []PatientLink    `json:"link,omitempty"`
}

// NewPatient creates a new Patient with the required fields
func NewPatient() *Patient {
	return &Patient{
		Base: Base{
			ResourceType: ResourceTypePatient,
		},
	}
}

// Address represents a physical address
type Address struct {
	Use        string   `json:"use,omitempty"`
	Type       string   `json:"type,omitempty"`
	Text       string   `json:"text,omitempty"`
	Line       []string `json:"line,omitempty"`
	City       string   `json:"city,omitempty"`
	District   string   `json:"district,omitempty"`
	State      string   `json:"state,omitempty"`
	PostalCode string   `json:"postalCode,omitempty"`
	Country    string   `json:"country,omitempty"`
	Period     *Period  `json:"period,omitempty"`
}

// Attachment represents a file or other attachment
type Attachment struct {
	ContentType string     `json:"contentType,omitempty"`
	Language    string     `json:"language,omitempty"`
	Data        string     `json:"data,omitempty"`
	URL         string     `json:"url,omitempty"`
	Size        *int       `json:"size,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	Title       string     `json:"title,omitempty"`
	Creation    *time.Time `json:"creation,omitempty"`
}

// PatientContact represents a patient's contact person
type PatientContact struct {
	Relationship []CodeableConcept `json:"relationship,omitempty"`
	Name         *HumanName        `json:"name,omitempty"`
	Telecom      []ContactPoint    `json:"telecom,omitempty"`
	Address      *Address          `json:"address,omitempty"`
	Gender       string            `json:"gender,omitempty"`
	Organization *Reference        `json:"organization,omitempty"`
	Period       *Period           `json:"period,omitempty"`
}

// Communication represents a patient's language preferences
type Communication struct {
	Language  CodeableConcept `json:"language"`
	Preferred *bool           `json:"preferred,omitempty"`
}

// PatientLink represents a link to another patient record
type PatientLink struct {
	Other Reference `json:"other"`
	Type  string    `json:"type"`
}

// Bundle represents a collection of resources
type Bundle struct {
	Base
	Type  string        `json:"type"`
	Total *int          `json:"total,omitempty"`
	Link  []BundleLink  `json:"link,omitempty"`
	Entry []BundleEntry `json:"entry,omitempty"`
}

// BundleLink represents a navigation link
type BundleLink struct {
	Relation string `json:"relation"`
	URL      string `json:"url"`
}

// BundleEntry represents a single entry in a bundle
type BundleEntry struct {
	FullURL  string          `json:"fullUrl,omitempty"`
	Resource json.RawMessage `json:"resource,omitempty"`
	Search   *BundleSearch   `json:"search,omitempty"`
}

// BundleSearch represents search information for a bundle entry
type BundleSearch struct {
	Mode  string   `json:"mode,omitempty"`
	Score *float64 `json:"score,omitempty"`
}

// GetTypedResource converts a raw resource to its typed struct using the resource mapper
func (b *Bundle) GetTypedResource(data json.RawMessage) (Resource, error) {
	if data == nil {
		return nil, nil
	}

	// Create a new resource mapper
	mapper := NewResourceMapper()
	return mapper.UnmarshalResource(data)
}

// UnmarshalJSON implements custom JSON unmarshaling for Patient
func (p *Patient) UnmarshalJSON(data []byte) error {
	type Alias Patient
	aux := struct {
		*Alias
		BirthDate  string `json:"birthDate,omitempty"`
		DeceasedAt string `json:"deceasedDateTime,omitempty"`
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.BirthDate != "" {
		t, err := time.Parse("2006-01-02", aux.BirthDate)
		if err != nil {
			return err
		}
		p.BirthDate = &t
	}

	if aux.DeceasedAt != "" {
		t, err := time.Parse(time.RFC3339, aux.DeceasedAt)
		if err != nil {
			return err
		}
		p.DeceasedAt = &t
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for Patient
func (p Patient) MarshalJSON() ([]byte, error) {
	type Alias Patient
	aux := struct {
		*Alias
		BirthDate  string `json:"birthDate,omitempty"`
		DeceasedAt string `json:"deceasedDateTime,omitempty"`
	}{
		Alias: (*Alias)(&p),
	}

	if p.BirthDate != nil {
		aux.BirthDate = p.BirthDate.Format("2006-01-02")
	}

	if p.DeceasedAt != nil {
		aux.DeceasedAt = p.DeceasedAt.Format(time.RFC3339)
	}

	return json.Marshal(aux)
}

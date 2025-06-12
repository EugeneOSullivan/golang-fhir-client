package search

import (
	"net/url"
	"strconv"
	"strings"
)

// Parameters represents FHIR search parameters
type Parameters struct {
	params url.Values
}

// NewParameters creates a new search parameters instance
func NewParameters() *Parameters {
	return &Parameters{
		params: make(url.Values),
	}
}

// Add adds a simple search parameter
func (p *Parameters) Add(name, value string) *Parameters {
	p.params.Add(name, value)
	return p
}

// AddModifier adds a search parameter with a modifier
func (p *Parameters) AddModifier(name, modifier, value string) *Parameters {
	p.params.Add(name+":"+modifier, value)
	return p
}

// AddPrefix adds a search parameter with a prefix
func (p *Parameters) AddPrefix(name, prefix, value string) *Parameters {
	p.params.Add(name, prefix+value)
	return p
}

// Count sets the _count parameter
func (p *Parameters) Count(count int) *Parameters {
	p.params.Set("_count", strconv.Itoa(count))
	return p
}

// Sort adds a sort parameter
func (p *Parameters) Sort(field string, desc bool) *Parameters {
	if desc {
		field = "-" + field
	}
	p.params.Add("_sort", field)
	return p
}

// Include adds an _include parameter
func (p *Parameters) Include(resourceType, searchParam string) *Parameters {
	p.params.Add("_include", resourceType+":"+searchParam)
	return p
}

// RevInclude adds a _revinclude parameter
func (p *Parameters) RevInclude(resourceType, searchParam string) *Parameters {
	p.params.Add("_revinclude", resourceType+":"+searchParam)
	return p
}

// Filter adds a filter parameter
func (p *Parameters) Filter(expression string) *Parameters {
	p.params.Add("_filter", expression)
	return p
}

// Elements sets the _elements parameter for field filtering
func (p *Parameters) Elements(fields ...string) *Parameters {
	p.params.Set("_elements", strings.Join(fields, ","))
	return p
}

// Summary sets the _summary parameter
func (p *Parameters) Summary(value string) *Parameters {
	p.params.Set("_summary", value)
	return p
}

// Total sets the _total parameter
func (p *Parameters) Total(value string) *Parameters {
	p.params.Set("_total", value)
	return p
}

// Format sets the _format parameter
func (p *Parameters) Format(format string) *Parameters {
	p.params.Set("_format", format)
	return p
}

// Since adds a _since parameter
func (p *Parameters) Since(timestamp string) *Parameters {
	p.params.Set("_since", timestamp)
	return p
}

// At adds an _at parameter
func (p *Parameters) At(timestamp string) *Parameters {
	p.params.Set("_at", timestamp)
	return p
}

// Profile adds a _profile parameter
func (p *Parameters) Profile(url string) *Parameters {
	p.params.Add("_profile", url)
	return p
}

// Security adds a _security parameter
func (p *Parameters) Security(system, code string) *Parameters {
	p.params.Add("_security", system+"|"+code)
	return p
}

// Tag adds a _tag parameter
func (p *Parameters) Tag(system, code string) *Parameters {
	p.params.Add("_tag", system+"|"+code)
	return p
}

// Contains adds a :contains modifier
func (p *Parameters) Contains(field, value string) *Parameters {
	return p.AddModifier(field, "contains", value)
}

// Exact adds an :exact modifier
func (p *Parameters) Exact(field, value string) *Parameters {
	return p.AddModifier(field, "exact", value)
}

// Missing adds a :missing modifier
func (p *Parameters) Missing(field string, isMissing bool) *Parameters {
	return p.AddModifier(field, "missing", strconv.FormatBool(isMissing))
}

// Type adds a :type modifier
func (p *Parameters) Type(field, resourceType string) *Parameters {
	return p.AddModifier(field, "type", resourceType)
}

// Above adds a gt (greater than) prefix
func (p *Parameters) Above(field, value string) *Parameters {
	return p.AddPrefix(field, "gt", value)
}

// Below adds a lt (less than) prefix
func (p *Parameters) Below(field, value string) *Parameters {
	return p.AddPrefix(field, "lt", value)
}

// EqualTo adds an eq prefix
func (p *Parameters) EqualTo(field, value string) *Parameters {
	return p.AddPrefix(field, "eq", value)
}

// NotEqualTo adds a ne prefix
func (p *Parameters) NotEqualTo(field, value string) *Parameters {
	return p.AddPrefix(field, "ne", value)
}

// GreaterThanOrEqual adds a ge prefix
func (p *Parameters) GreaterThanOrEqual(field, value string) *Parameters {
	return p.AddPrefix(field, "ge", value)
}

// LessThanOrEqual adds a le prefix
func (p *Parameters) LessThanOrEqual(field, value string) *Parameters {
	return p.AddPrefix(field, "le", value)
}

// StartsWith adds a sw prefix
func (p *Parameters) StartsWith(field, value string) *Parameters {
	return p.AddPrefix(field, "sw", value)
}

// EndsWith adds an ew prefix
func (p *Parameters) EndsWith(field, value string) *Parameters {
	return p.AddPrefix(field, "ew", value)
}

// Encode returns the encoded URL query string
func (p *Parameters) Encode() string {
	return p.params.Encode()
}

// Raw returns the raw parameter map
func (p *Parameters) Raw() url.Values {
	return p.params
}

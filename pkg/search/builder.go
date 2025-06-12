package search

import (
	"net/url"
	"strconv"
	"strings"
)

// Builder represents a FHIR search query builder
type Builder struct {
	params url.Values
}

// NewBuilder creates a new search builder
func NewBuilder() *Builder {
	return &Builder{
		params: make(url.Values),
	}
}

// Where adds a search parameter with an operator and value
func (b *Builder) Where(param, op, value string) *Builder {
	key := param
	if op != "" && op != "eq" {
		key = param + ":" + op
	}
	b.params.Add(key, value)
	return b
}

// Sort adds a sort parameter
func (b *Builder) Sort(param string, desc bool) *Builder {
	if desc {
		param = "-" + param
	}
	b.params.Add("_sort", param)
	return b
}

// Count sets the _count parameter
func (b *Builder) Count(count int) *Builder {
	b.params.Set("_count", strconv.Itoa(count))
	return b
}

// Include adds an _include parameter
func (b *Builder) Include(resourceType, searchParam string) *Builder {
	b.params.Add("_include", resourceType+":"+searchParam)
	return b
}

// RevInclude adds a _revinclude parameter
func (b *Builder) RevInclude(resourceType, searchParam string) *Builder {
	b.params.Add("_revinclude", resourceType+":"+searchParam)
	return b
}

// Build returns the encoded URL query string
func (b *Builder) Build() string {
	return b.params.Encode()
}

// BuildRaw returns the unencoded URL query string
func (b *Builder) BuildRaw() string {
	var pairs []string
	for key, values := range b.params {
		for _, value := range values {
			pairs = append(pairs, key+"="+value)
		}
	}
	return strings.Join(pairs, "&")
}

package glue

import (
	"strings"
)

//go:generate stringer -type=MapKeyFilterType map_filter_type.go

// MapKeyFilterType is an enum of the type of map key filter.
type MapKeyFilterType int

const (
	// TypeKeyEquals ...
	TypeKeyEquals MapKeyFilterType = iota
	// TypeKeyContains ...
	TypeKeyContains
	// TypeKeyPrefix ...
	TypeKeyPrefix
	// TypeKeySuffix ...
	TypeKeySuffix
)

// MapKeyFilterFunc ...
type MapKeyFilterFunc func(key string, input interface{}) bool

type KeyFilter struct {
	Input interface{}
	Func  MapKeyFilterFunc
}

// MapFilter ...
type MapFilter struct {
	Input  map[string]interface{}
	Output map[string]interface{}

	keyFilters map[MapKeyFilterType]*KeyFilter
}

func (m *MapFilter) ApplyKeyFilter(
	filterType MapKeyFilterType,
	key string,
	value interface{}) bool {

	filterInput, keyFilterFunc, ok := m.KeyFilter(filterType)
	if !ok {
		return false
	}

	switch filterInput.(type) {
	case []interface{}:
		for _, filterV := range filterInput.([]interface{}) {
			if keyFilterFunc(key, filterV.(string)) {
				m.Output[key] = value
				return true
			}
		}
	case map[interface{}]bool:
		if keyFilterFunc(key, filterInput) {
			m.Output[key] = value
			return true
		}
	}

	return false
}

func (m *MapFilter) SetKeyFilter(name string, values []interface{}) bool {
	var filterType MapKeyFilterType
	var filterFunc MapKeyFilterFunc
	var filterInput interface{} = values

	switch name {
	case "equals":
		filterType = TypeKeyEquals
		filterInput = ArrayToSet(values)
		filterFunc = func(key string, input interface{}) bool {
			if set, ok := input.(map[interface{}]bool); ok {
				if _, ok := set[key]; ok {
					return true
				}
			}

			return false
		}
	case "contains":
		filterType = TypeKeyContains
		filterFunc = func(key string, input interface{}) bool {
			return strings.Contains(key, input.(string))
		}
	case "prefix":
		filterType = TypeKeyPrefix
		filterFunc = func(key string, input interface{}) bool {
			return strings.HasPrefix(key, input.(string))
		}
	case "suffix":
		filterType = TypeKeySuffix
		filterFunc = func(key string, input interface{}) bool {
			return strings.HasSuffix(key, input.(string))
		}
	default:
		return false
	}

	if m.keyFilters == nil {
		m.keyFilters = make(map[MapKeyFilterType]*KeyFilter)
	}
	m.keyFilters[filterType] = &KeyFilter{Input: filterInput, Func: filterFunc}
	return true
}

func (m *MapFilter) KeyFilter(filterType MapKeyFilterType) (interface{}, MapKeyFilterFunc, bool) {
	if filter, ok := m.keyFilters[filterType]; ok {
		return filter.Input, filter.Func, true
	}

	return nil, nil, false
}

// Apply ..
func (m *MapFilter) Apply() {
	m.Output = make(map[string]interface{})

	for k, v := range m.Input {
		// Skip this key if already previously matched by a filter
		if _, ok := m.Output[k]; ok {
			continue
		}

		if m.ApplyKeyFilter(TypeKeyEquals, k, v) {
			continue
		}

		if m.ApplyKeyFilter(TypeKeyContains, k, v) {
			continue
		}

		if m.ApplyKeyFilter(TypeKeyPrefix, k, v) {
			continue
		}

		m.ApplyKeyFilter(TypeKeySuffix, k, v)
	}

	// Set the output to the unfiltered input if no keys matched by any filter
	if len(m.Output) == 0 {
		m.Output = m.Input
	}
}

package filters

import (
	"fmt"
	"strings"
)

type Filters []Filter

func (f *Filters) GetFilters() []Filter {
	return *f
}

func (f *Filters) GetFilter(key FilterKey) Filter {
	for _, filter := range *f {
		if filter.FilterType() == FilterTypeKeyValue && filter.(*FilterKeyValue).Key == key {
			return filter
		}
	}
	return nil
}

func (f *Filters) AddFilter(filter Filter) {
	*f = append(*f, filter)
}

func (f *Filters) GetLabelFilter() *LabelFilter {
	for _, filter := range *f {
		if labelFilter, ok := filter.(*LabelFilter); ok {
			return labelFilter
		}
	}
	return nil
}

func (f *Filters) GetKeyValueFilter(key FilterKey) *FilterKeyValue {
	for _, filter := range *f {
		if keyValueFilter, ok := filter.(*FilterKeyValue); ok {
			if strings.EqualFold(string(keyValueFilter.Key), string(key)) {
				return keyValueFilter
			}
		}
	}
	return nil
}

type Filter interface {
	FilterType() FilterType
	ToParams() map[string]string
}

type FilterType string

const (
	FilterTypeKeyValue FilterType = "keyValue"
	FilterTypeLabel    FilterType = "label"
)

type FilterKeyValue struct {
	Key   FilterKey `json:"key"`
	Value string    `json:"value"`
}

func (f *FilterKeyValue) FilterType() FilterType {
	return FilterTypeKeyValue
}

type FilterKey string

const (
	FilterRegion               FilterKey = "region"
	FilterZone                 FilterKey = "zone"
	FilterVpcIdentity          FilterKey = "vpc"
	FilterSubnetIdentity       FilterKey = "subnet"
	FilterMachineIdentity      FilterKey = "machine"
	FilterLoadbalancerIdentity FilterKey = "loadbalancer"
	FilterReservedIp           FilterKey = "reservedIp"
	FilterName                 FilterKey = "name"
	FilterIdentity             FilterKey = "identity"
	FilterSlug                 FilterKey = "slug"
	FilterStatus               FilterKey = "status"
)

type LabelFilter struct {
	MatchLabels map[string]string `json:"matchLabels"`
}

func (f *LabelFilter) FilterType() FilterType {
	return FilterTypeLabel
}

// Parses to query params like matchLabels[env]=prod from map["env"] = "prod"
func (f *LabelFilter) ToParams() map[string]string {
	params := map[string]string{}
	for k, v := range f.MatchLabels {
		params[fmt.Sprintf("matchLabels[%s]", k)] = v
	}
	return params
}

func (f *FilterKeyValue) ToParams() map[string]string {
	if strings.TrimSpace(string(f.Key)) == "" {
		return map[string]string{}
	}
	if strings.TrimSpace(f.Value) == "" {
		return map[string]string{}
	}
	return map[string]string{
		string(f.Key): f.Value,
	}
}

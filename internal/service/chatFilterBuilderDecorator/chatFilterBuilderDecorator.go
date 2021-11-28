package chatFilterBuilderDecorator

import (
	"strings"

	"github.com/darthbanana13/datastorageapi/internal/filterBuilder"
)

type FilterParams struct {
	filter filterBuilder.FilterParams
}

func NewChatFilter(filter filterBuilder.FilterParams) filterBuilder.FilterParams {
	return &FilterParams{filter: filter}
}

func (f *FilterParams) GetFieldConditions() map[string]interface{} {
	fields := f.filter.GetFieldConditions()
	if val, ok := fields["language"]; ok {
		fields["language"] = strings.ToLower(val.(string))
	}
	return fields
}

func (f *FilterParams) GetOffset() uint {
	return f.filter.GetOffset()
}

func (f *FilterParams) GetLimit() uint {
	return f.filter.GetLimit()
}

func (f *FilterParams) GetSortFields() map[string]string {
	return f.filter.GetSortFields()
}

func (f *FilterParams) GetReturnFields() []string {
	return f.filter.GetReturnFields()
}

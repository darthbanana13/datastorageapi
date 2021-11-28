package filterBuilder

//TODO: Should probably get rid of internal dependencies & move it to pkg
import (
	"github.com/darthbanana13/datastorageapi/internal/aqlBuilder"
)

type Filter struct {
	conditions   map[string]interface{}
	offset       uint
	limit        uint
	sort         map[string]string
	returnFields []string
}

//TODO: This interface is a tad on the fat side, consider slimming it down
type FilterParams interface {
	GetFieldConditions() map[string]interface{}
	GetOffset() uint
	GetLimit() uint
	GetSortFields() map[string]string
	GetReturnFields() []string
}

func NewFilter() Filter {
	return Filter{
		conditions: make(map[string]interface{}),
		sort:       make(map[string]string),
	}
}

func (f *Filter) WithFieldCodition(fieldName string, fieldValue interface{}) {
	f.conditions[fieldName] = fieldValue
}

func (f *Filter) WithOffsetAndLimit(offset, limit uint) {
	f.offset = offset
	f.limit = limit
}

func (f *Filter) WithSortFieldDescending(fieldName string) {
	f.sort[fieldName] = aqlBuilder.Descending
}

func (f *Filter) WithReturnField(fieldName string) {
	f.returnFields = append(f.returnFields, fieldName)
}

func (f *Filter) GetFieldConditions() map[string]interface{} {
	return f.conditions
}

func (f *Filter) GetOffset() uint {
	return f.offset
}

func (f *Filter) GetLimit() uint {
	return f.limit
}

func (f *Filter) GetSortFields() map[string]string {
	return f.sort
}

func (f *Filter) GetReturnFields() []string {
	return f.returnFields
}

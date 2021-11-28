package chatFilterBuilderDecorator_test

import (
	"testing"

	"github.com/darthbanana13/datastorageapi/internal/filterBuilder"
	filterDecorator "github.com/darthbanana13/datastorageapi/internal/service/chatFilterBuilderDecorator"

	"github.com/stretchr/testify/assert"
)

func TestLanguageToLower(t *testing.T) {
	filterBuilder := filterBuilder.NewFilter()
	filterBuilder.WithFieldCodition("language", "DE")
	f := filterDecorator.NewChatFilter(&filterBuilder)
	assert.Equal(t, f.GetFieldConditions(), map[string]interface{}{"language": "de"})
}

func TestLanguageToLowerOnLowercase(t *testing.T) {
	filterBuilder := filterBuilder.NewFilter()
	filterBuilder.WithFieldCodition("language", "it")
	filterBuilder.WithFieldCodition("nonImportantField", "Other VALUE")
	f := filterDecorator.NewChatFilter(&filterBuilder)
	assert.Equal(
		t,
		f.GetFieldConditions(),
		map[string]interface{}{"language": "it", "nonImportantField": "Other VALUE"},
	)
}

func TestLanguageToLowerOnOtherFields(t *testing.T) {
	filterBuilder := filterBuilder.NewFilter()
	filterBuilder.WithFieldCodition("nonLanguageField", "Totally not a language")
	filterBuilder.WithFieldCodition("i18n", "素晴らしい")
	f := filterDecorator.NewChatFilter(&filterBuilder)
	assert.Equal(
		t,
		f.GetFieldConditions(),
		map[string]interface{}{
			"nonLanguageField": "Totally not a language",
			"i18n":             "素晴らしい",
		},
	)
}

package aqlBuilder_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/darthbanana13/datastorageapi/internal/aqlBuilder"
	"github.com/stretchr/testify/assert"
)

func whiteSpaceStripper(s string) string {
	space := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(space.ReplaceAllString(s, " "))
}

func TestInsert(t *testing.T) {
	builder := aqlBuilder.NewBuilder("documents")
	builder.WithInsert(map[string]interface{}{
		"field1":       "value1",
		"anotherField": 7,
	})
	actual := whiteSpaceStripper(builder.Build())

	expected := whiteSpaceStripper(`
				INSERT {
					Field1: @field1,
					AnotherField: @anotherField
				} INTO documents
	`)
	assert.Equal(t, expected, actual)
}

func TestUpdate(t *testing.T) {
	builder := aqlBuilder.NewBuilder("updater")
	builder.WithAndFilterCondition(map[string]interface{}{
		"someField":      5.75,
		"whatIsInTheBox": "bananas",
	})
	builder.WithLoopStatement()
	builder.WithUpdate(map[string]interface{}{
		"batSignalOn": true,
	})
	actual := whiteSpaceStripper(builder.Build())

	expected := whiteSpaceStripper(`
				FOR u IN updater
					FILTER u.SomeField == @someField
					FILTER u.WhatIsInTheBox == @whatIsInTheBox
					UPDATE u._key WITH {
						BatSignalOn: @batSignalOn
					} IN updater
	`)
	assert.Equal(t, expected, actual)
}

func TestRemove(t *testing.T) {
	builder := aqlBuilder.NewBuilder("remover")
	builder.WithAndFilterCondition(map[string]interface{}{
		"fileLocation": "some/compromising/file",
	})
	builder.WithRemove()
	builder.WithLoopStatement()
	actual := whiteSpaceStripper(builder.Build())

	expected := whiteSpaceStripper(`
				FOR r IN remover
					FILTER r.FileLocation == @fileLocation
					REMOVE r IN remover
	`)
	assert.Equal(t, expected, actual)
}

func TestFilter(t *testing.T) {
	builder := aqlBuilder.NewBuilder("enormousCollection")
	builder.WithLoopStatement()
	builder.WithAndFilterCondition(map[string]interface{}{
		"blockbusterLaunchYear": 2021,
		"genre":                 "Sci-Fi",
	})
	builder.WithSortFields(map[string]string{
		"imdbRating": aqlBuilder.Descending,
	})
	builder.WithLimit(5, 10)
	builder.WithReturnFields([]string{"title", "imdbRating"})
	actual := whiteSpaceStripper(builder.Build())

	expected := whiteSpaceStripper(`
				FOR e IN enormousCollection
					FILTER e.BlockbusterLaunchYear == @blockbusterLaunchYear
					FILTER e.Genre == @genre
					SORT e.ImdbRating DESC
					LIMIT @offset, @limit
					RETURN {
						Title: e.Title,
						ImdbRating: e.ImdbRating
					}
	`)

	assert.Equal(t, expected, actual)
}

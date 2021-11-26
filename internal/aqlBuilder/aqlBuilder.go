package aqlBuilder

import (
    "fmt"
    "strings"

	log "github.com/sirupsen/logrus"
)

type Builder struct {
    collectionName      string
    collectionInterator string
    loopStatement       string
    filterCondition     string
    sortFields          string
    returnStatement     string
    limitStatement      string
}

const Ascending = "ASC"
const Descending = "DESC"

func NewBuilder(collectionName string) Builder {
    return Builder{
        collectionName: collectionName,
        collectionInterator: string(collectionName[0]),
    }
}

func (b *Builder) WithLoopStatement() {
    b.loopStatement += "FOR " + b.collectionInterator + " IN " + b.collectionName + "\n"
}

//I know, fields should be a []string but the conversions would take useless CPU cycles IMO
func (b *Builder) WithAndFilterCondition(fields map[string]interface{}) {
    for k := range fields {
        b.filterCondition += fmt.Sprintf(
            "FILTER %s.%s == @%s\n",
            b.collectionInterator,
            strings.Title(k),
            k,
        )
    }
}

func (b *Builder) WithSortCondition(fields map[string]string) {
    var sortFields []string
    for k, v := range fields {
        sortFields = append(
            sortFields,
            b.collectionInterator + "." + strings.Title(k) + " " + v,
        )
    }
    b.sortFields += "SORT " + strings.Join(sortFields, ", ") + "\n"
}

func (b *Builder) WithReturnFields(fields []string) {
    var returnFields []string
    for _, v := range fields {
        returnFields = append(
            returnFields,
            strings.Title(v) + ": " + b.collectionInterator + "." + strings.Title(v),
        )
    }
    b.returnStatement += "RETURN {\n" + strings.Join(returnFields, ",\n") + "\n}\n"
}

func (b *Builder) WithLimit() {
    b.limitStatement += "LIMIT @offset, @limit\n"
}

func (b *Builder) Build() string {
    return b.loopStatement + b.filterCondition + b.sortFields + b.limitStatement + b.returnStatement
}

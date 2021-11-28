package aqlBuilder

import (
	"strings"

	driver "github.com/arangodb/go-driver"
	container "github.com/golobby/container/v3"
	log "github.com/sirupsen/logrus"
)

type Builder struct {
	collectionName      string
	collectionInterator string
	loopStatement       string
	insertStatement     string
	filterCondition     string
	removeStatement     string
	sortFields          string
	updateStatement     string
	returnStatement     string
	limitStatement      string
	boundFields         map[string]interface{}
}

const Ascending = "ASC"
const Descending = "DESC"

func NewBuilder(collectionName string) Builder {
	return Builder{
		collectionName:      collectionName,
		collectionInterator: string(collectionName[0]),
		boundFields:         make(map[string]interface{}),
	}
}

func (b *Builder) WithInsert(fields map[string]interface{}) {
	var insertFields []string
	for k, v := range fields {
		insertFields = append(
			insertFields,
			strings.Title(k)+": "+"@"+k,
		)
		b.boundFields[k] = v
	}
	b.insertStatement += "INSERT\n{\n" + strings.Join(insertFields, ",\n") + "\n}\nINTO " + b.collectionName + "\n"
}

func (b *Builder) WithLoopStatement() {
	b.loopStatement += "FOR " + b.collectionInterator + " IN " + b.collectionName + "\n"
}

func (b *Builder) WithAndFilterCondition(fields map[string]interface{}) {
	var filterConditions []string
	for k, v := range fields {
		filterConditions = append(
			filterConditions,
			"FILTER "+b.collectionInterator+"."+strings.Title(k)+" == @"+k,
		)
		b.boundFields[k] = v
	}
	b.filterCondition += strings.Join(filterConditions, "\n") + "\n"
}

func (b *Builder) WithSortCondition(fields map[string]string) {
	var sortFields []string
	for k, v := range fields {
		sortFields = append(
			sortFields,
			b.collectionInterator+"."+strings.Title(k)+" "+v,
		)
	}
	b.sortFields += "SORT " + strings.Join(sortFields, ", ") + "\n"
}

func (b *Builder) WithReturnFields(fields []string) {
	var returnFields []string
	for _, v := range fields {
		returnFields = append(
			returnFields,
			strings.Title(v)+": "+b.collectionInterator+"."+strings.Title(v),
		)
	}
	b.returnStatement += "RETURN {\n" + strings.Join(returnFields, ",\n") + "\n}\n"
}

func (b *Builder) WithLimit(offset, limit uint) {
	b.limitStatement += "LIMIT @offset, @limit\n"
	b.boundFields["offset"] = offset
	b.boundFields["limit"] = limit
}

func (b *Builder) WithRemove() {
	b.removeStatement += "REMOVE " + b.collectionInterator + " IN " + b.collectionName + "\n"
}

func (b *Builder) WithUpdate(fields map[string]interface{}) {
	var updateFields []string
	for k, v := range fields {
		updateFields = append(
			updateFields,
			strings.Title(k)+": "+"@"+k,
		)
		b.boundFields[k] = v
	}
	b.updateStatement += "UPDATE " + b.collectionInterator + "._key WITH {" + strings.Join(updateFields, ",\n") + "\n} IN " + b.collectionName + "\n"
}

func (b *Builder) Execute() (driver.Cursor, error) {
	db := loadDbConn()
	query := b.Build()
	cursor, err := (*db).Query(nil, query, b.boundFields)
	if err != nil {
		log.Errorf("Upsy! This query has a problem:\n%s\n%s", query, err)
	}
	return cursor, err
}

//TODO: Should probably do something about each With* function adding a new line at the end then trimming the remaining new lines
func (b *Builder) Build() string {
	retrunString := b.loopStatement
	retrunString += b.insertStatement
	retrunString += b.filterCondition
	retrunString += b.updateStatement
	retrunString += b.sortFields
	retrunString += b.limitStatement
	retrunString += b.returnStatement
	retrunString += b.removeStatement
	return strings.TrimSpace(retrunString)
}

//TODO: Should be moved to another place or it's probably not that KISS if I just use the container here
func loadDbConn() *driver.Database {
	var db *driver.Database
	err := container.Resolve(&db)
	if err != nil {
		log.Panic("Unable to get DB connection")
	}
	return db
}

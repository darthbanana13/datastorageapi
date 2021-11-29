package cursorIterator

import (
    driver "github.com/arangodb/go-driver"

	log "github.com/sirupsen/logrus"
)

func ToMap(cursor driver.Cursor) ([]map[string]interface{}, error) {
	defer cursor.Close()
	var docs []map[string]interface{}

	for {
		var doc map[string]interface{}
        _, err := cursor.ReadDocument(nil, &doc)

		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Errorf("Got this document from the DB: %v", err)
			return []map[string]interface{}{}, err
		} else {
			docs = append(docs, doc)
		}
	}
	return docs, nil
}

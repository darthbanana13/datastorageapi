package chat

import (
	"strings"

	aql "github.com/darthbanana13/datastorageapi/internal/aqlBuilder"
	"github.com/darthbanana13/datastorageapi/internal/model/chat"
	"github.com/darthbanana13/datastorageapi/internal/cursorIterator"

	driver "github.com/arangodb/go-driver"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func Insert(msg chat.ChatMsg) error {
	var mapChat map[string]interface{}
	err := mapstructure.Decode(msg, &mapChat)
	if err != nil {
		log.Error("Could not convert struct to map", err)
		return err
	}

	aqlBuilder := aql.NewBuilder(chat.CollName)
	aqlBuilder.WithInsert(mapChat)
	_, err = aqlBuilder.Execute()
	return err
}

func DeleteDialog(dialogId uint) error {
	aqlBuilder := aql.NewBuilder(chat.CollName)
	aqlBuilder.WithLoopStatement()
	aqlBuilder.WithAndFilterCondition(map[string]interface{}{"dialogId": dialogId})
	aqlBuilder.WithRemove()
	_, err := aqlBuilder.Execute()
	return err
}

//TODO: This function isn't used anywhere yet
func DeleteCustomer(customerId uint) error {
	aqlBuilder := aql.NewBuilder(chat.CollName)
	aqlBuilder.WithLoopStatement()
	aqlBuilder.WithAndFilterCondition(map[string]interface{}{"customerId": customerId})
	aqlBuilder.WithRemove()
	_, err := aqlBuilder.Execute()
	return err
}

func ConsentTrueDialog(dialogId uint) error {
	aqlBuilder := aql.NewBuilder(chat.CollName)
	aqlBuilder.WithLoopStatement()
	aqlBuilder.WithAndFilterCondition(map[string]interface{}{"dialogId": dialogId})
	aqlBuilder.WithUpdate(map[string]interface{}{"consent": true})
	_, err := aqlBuilder.Execute()
	return err
}

//TODO: Find a clean expressive way to add a ton of parameters
func AndFilter(fieldConditions map[string]interface{}, offset, limit uint) ([]map[string]interface{}, error) {
	if val, ok := fieldConditions["language"]; ok {
		fieldConditions["language"] = strings.ToLower(val.(string))
	}

	aqlBuilder := aql.NewBuilder(chat.CollName)
	aqlBuilder.WithLoopStatement()
	aqlBuilder.WithAndFilterCondition(fieldConditions)
	aqlBuilder.WithSortCondition(map[string]string{"NanoTimestamp": aql.Descending})
	aqlBuilder.WithReturnFields([]string{"Text", "Language", "CustomerId", "DialogId"})
	aqlBuilder.WithLimit(offset, limit)
	cursor, err := aqlBuilder.Execute()

	if err != nil {
		return []map[string]interface{}{}, err
	}
	return cursoriterator.ToMap(cursor)
}

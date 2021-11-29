package chat

import (
	aql "github.com/darthbanana13/datastorageapi/internal/aqlBuilder"
	"github.com/darthbanana13/datastorageapi/internal/cursorIterator"
	"github.com/darthbanana13/datastorageapi/internal/filterBuilder"
	"github.com/darthbanana13/datastorageapi/internal/model/chat"

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

func AndFilter(filterParams filterBuilder.FilterParams) ([]map[string]interface{}, error) {
	aqlBuilder := aql.NewBuilder(chat.CollName)
	aqlBuilder.WithLoopStatement()
	aqlBuilder.WithAndFilterCondition(filterParams.GetFieldConditions())
	aqlBuilder.WithSortFields(filterParams.GetSortFields())
	aqlBuilder.WithReturnFields(filterParams.GetReturnFields())
	aqlBuilder.WithLimit(filterParams.GetOffset(), filterParams.GetLimit())
	cursor, err := aqlBuilder.Execute()

	if err != nil {
		return []map[string]interface{}{}, err
	}
	return cursorIterator.ToMap(cursor)
}

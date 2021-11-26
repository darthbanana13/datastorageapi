package chat

import (
	"time"
	"strings"
	"encoding/json"

	aql "github.com/darthbanana13/datastorageapi/internal/aqlBuilder"

	driver "github.com/arangodb/go-driver"
	container "github.com/golobby/container/v3"
	log "github.com/sirupsen/logrus"
)

const collName = "chat"

type ChatMsg struct {
	Text          string
	Language      string
	CustomerId    uint
	DialogId      uint
	Consent       bool
	NanoTimestamp int64
}

func NewChatMsg(customerId, dialogId uint, text, language string) ChatMsg {
	var msg ChatMsg

	msg.CustomerId = customerId
	msg.DialogId = dialogId
	msg.Text = text
	msg.Language = strings.ToLower(language)
	msg.Consent = false
	msg.NanoTimestamp = time.Now().UnixNano()

	return msg
}

//TODO: Should be moved to another place
//TODO: Maybe it isn't the best idea to use this function everywhere, probably have db as a global package variable?
func loadDbConn() *driver.Database {
	var db *driver.Database
	err := container.Resolve(&db)
	if err != nil {
		log.Panic("Unable to get DB connection")
	}
	return db
}

func (msg *ChatMsg) Insert() error {
	db := loadDbConn()
	byteChat, err := json.Marshal(msg)
	if err != nil {
		log.Error("Unable to marshal chatMsg data", err)
		return err
	}
	_, err = (*db).Query(nil, "INSERT "+string(byteChat)+" INTO "+collName, nil)
	return err
}

func DeleteDialog(dialogId uint) error {
	db := loadDbConn()
	_, err := (*db).Query(
		nil,
		`FOR c IN chat
			FILTER c.DialogId == @dialogId
			REMOVE c IN chat`,
		map[string]interface{}{"dialogId": dialogId},
	)
	return err
}

//TODO: This function isn't used anywhere yet
func DeleteCustomer(customerId uint) error {
	db := loadDbConn()
	_, err := (*db).Query(
		nil,
		`FOR c IN chat
			FILTER c.CustomerId == @customerId
			REMOVE c IN chat`,
		map[string]interface{}{"customerId": customerId},
	)
	return err
}

func ConsentTrueDialog(dialogId uint) error {
	db := loadDbConn()
	_, err := (*db).Query(
		nil,
		`FOR c IN chat
			FILTER c.DialogId == @dialogId
			UPDATE c._key WITH { Consent: true } IN chat`,
		map[string]interface{}{"dialogId": dialogId},
	)
	return err
}

//TODO: Refactor this BS
//TODO: Find a clean expressive way to add a ton of parameters
func AndFilter(fieldConditions map[string]interface{}, offset, limit uint) ([]map[string]interface{}, error) {
	db := loadDbConn()

	aqlBuilder := aql.NewBuilder(collName)
	aqlBuilder.WithLoopStatement()
	aqlBuilder.WithAndFilterCondition(fieldConditions)
	aqlBuilder.WithSortCondition(map[string]string{"NanoTimestamp": aql.Descending})
	aqlBuilder.WithReturnFields([]string{"Text", "Language", "CustomerId", "DialogId"})
	aqlBuilder.WithLimit()

	fieldConditions["offset"] = offset
	fieldConditions["limit"] = limit

	query := aqlBuilder.Build()

	cursor, err := (*db).Query(nil, query, fieldConditions)
	if err != nil {
		log.Errorf("Upsy! This query has a problem:\n%s\n%s", query, err)
		return []map[string]interface{}{}, err
	}
	defer cursor.Close()

	var msgs []map[string]interface{}

	for {
		var msg map[string]interface{}
		_, err = cursor.ReadDocument(nil, &msg)

		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Errorf("Got this chat from the DB: %v", err)
			return []map[string]interface{}{}, err
		} else {
			msgs = append(msgs, msg)
		}
	}

	return msgs, nil
}

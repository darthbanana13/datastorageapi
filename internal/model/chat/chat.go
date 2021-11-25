package chat

import (
	"encoding/json"
	"time"

	driver "github.com/arangodb/go-driver"
	container "github.com/golobby/container/v3"
	log "github.com/sirupsen/logrus"
)

const collName = "chat"

type chatMsg struct {
	Text          string
	Language      string
	CustomerId    uint
	DialogId      uint
	Consent       bool
	NanoTimestamp int64
}

func NewChatMsg(customerId, dialogId uint, text, language string) chatMsg {
	var msg chatMsg

	msg.CustomerId = customerId
	msg.DialogId = dialogId
	msg.Consent = false
	msg.NanoTimestamp = time.Now().UnixNano()

	return msg
}

//TODO: Should be moved to another place
func loadDbConn() *driver.Database {
	var db *driver.Database
	err := container.Resolve(&db)
	if err != nil {
		log.Panic("Unable to get DB connection")
	}
	return db
}

func (msg *chatMsg) Insert() error {
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
	log.Info("DeleteDialog")
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
	log.Info("ConsentTrueDialog")
	_, err := (*db).Query(
		nil,
		`FOR c IN chat
			FILTER c.DialogId == @dialogId
			UPDATE c._key WITH { Consent: true } IN chat`,
		map[string]interface{}{"dialogId": dialogId},
	)
	return err
}

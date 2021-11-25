package chat

import (
	"encoding/json"
	"time"

	driver "github.com/arangodb/go-driver"
	container "github.com/golobby/container/v3"
	log "github.com/sirupsen/logrus"
)

type chatMsg struct {
	Text          string
	Language      string
	CustomerId    uint
	DialogId      uint
	Consent       bool
	NanoTimestamp int64
}

//NOTE: This isn't used anywhere, but seemed appropiate to have it
func NewChatMsg(customerId, dialogId uint, text, language string) chatMsg {
	var msg chatMsg

	msg.CustomerId = customerId
	msg.DialogId = dialogId
	msg.Consent = false
	msg.NanoTimestamp = time.Now().UnixNano()

	return msg
}

//TODO: Should be moved to another place
func lazyLoadDbConn() *driver.Database {
	var db *driver.Database
	err := container.Resolve(&db)
	if err != nil {
		log.Panic("Unable to get DB connection")
	}
	return db
}

func (msg *chatMsg) Insert() {
	db := lazyLoadDbConn()
	byteChat, err := json.Marshal(msg)
	if err != nil {
		log.Error("Unable to marshal data")
	}
	(*db).Query(nil, "INSERT"+string(byteChat)+"INTO chat", nil)
}

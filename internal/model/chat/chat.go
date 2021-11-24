package chat

import (
	"time"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	driver "github.com/arangodb/go-driver"
)

type chatMsg struct {
  Text		string `json:"text"`
  Language	string `json:"language"`
  CustomerId	int
  DialogId	int
  Consent	bool
  NanoTimestamp int64
}

//NOTE: This isn't used anywhere, but seemed appropiate to have it
// func NewChatMsg(customerId, dialogId int, text, language string) chatMsg {
  // var msg chatMsg
  // msg.CustomerId = customerId
  // msg.DialogId = dialogId
  // msg.Consent = false
  // msg.NanoTimestamp = time.Now().UnixNano()

  // return msg
// }

func NewEmptyChatMsg(customerId, dialogId int) chatMsg {
  var msg chatMsg
  msg.CustomerId = customerId
  msg.DialogId = dialogId
  msg.Consent = false
  msg.NanoTimestamp = time.Now().UnixNano()
  return msg
}

//TODO: should not have to pass the database every time
func (msg *chatMsg) Insert(db driver.Database) {
  byteChat, err := json.Marshal(msg)
  if err != nil {
    log.Error("Unable to marshal data")
  }
  db.Query(nil, "INSERT" + string(byteChat) + "INTO chat", nil)
}

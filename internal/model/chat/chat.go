package chat

import (
	"strings"
	"time"
)

const CollName = "chat"

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

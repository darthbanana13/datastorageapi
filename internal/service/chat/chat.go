package chat

import (
	"github.com/darthbanana13/datastorageapi/internal/model/chat"
)

func SaveMsg(customerId, dialogId uint, text, language string) error {
	msg := chat.NewChatMsg(
		customerId,
		dialogId,
		text,
		language,
	)
	return msg.Insert()
}

func Consent(dialogId uint, consent bool) error {
	if consent {
		return chat.ConsentTrueDialog(dialogId)
	}
	return chat.DeleteDialog(dialogId)
}

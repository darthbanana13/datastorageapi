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

func Filter(allFields map[string]interface{}, page, entriesPerPage uint) ([]map[string]interface{}, error) {
	fields := make(map[string]interface{})
	for k, v := range allFields {
		if s, ok := v.(string); ok == true && s != "" {
			fields[k] = v
		} else if i, ok := v.(uint); ok == true && i != 0 {
			fields[k] = v
		} else if _, ok := v.(bool); ok == true {
			fields[k] = v
		}
	}
	page = page - 1
	return chat.AndFilter(fields, entriesPerPage*page, entriesPerPage)
}

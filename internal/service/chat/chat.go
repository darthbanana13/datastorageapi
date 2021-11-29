package chat

import (
	"github.com/darthbanana13/datastorageapi/internal/filterBuilder"
	"github.com/darthbanana13/datastorageapi/internal/model/chat"
	chatRepository "github.com/darthbanana13/datastorageapi/internal/repository/chat"
	filterDecorator "github.com/darthbanana13/datastorageapi/internal/service/chatFilterBuilderDecorator"
)

func SaveMsg(customerId, dialogId uint, text, language string) error {
	msg := chat.NewChatMsg(
		customerId,
		dialogId,
		text,
		language,
	)
	return chatRepository.Insert(msg)
}

func Consent(dialogId uint, consent bool) error {
	if consent {
		return chatRepository.ConsentTrueDialog(dialogId)
	}
	return chatRepository.DeleteDialog(dialogId)
}

//TODO: Make this function smaller
func FilterAndSortByNewFirst(allFields map[string]interface{}, page, entriesPerPage uint) ([]map[string]interface{}, error) {
	filterParams := filterBuilder.NewFilter()
	for k, v := range allFields {
		//TODO: These or conditions could probably be written prettier
		if s, ok := v.(string); ok == true && s != "" {
			filterParams.WithFieldCodition(k, v)
		} else if i, ok := v.(uint); ok == true && i != 0 {
			filterParams.WithFieldCodition(k, v)
		} else if _, ok := v.(bool); ok == true {
			filterParams.WithFieldCodition(k, v)
		}
	}

	filterParams.WithOffsetAndLimit(entriesPerPage*(page-1), entriesPerPage)
	filterParams.WithSortFieldDescending("NanoTimestamp")
	filterParams.WithReturnField("Text")
	filterParams.WithReturnField("Language")
	filterParams.WithReturnField("CustomerId")
	filterParams.WithReturnField("DialogId")

	return chatRepository.AndFilter(filterDecorator.NewChatFilter(&filterParams))
}

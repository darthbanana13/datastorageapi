package chat

import (
	"encoding/json"
	"net/http"

	"github.com/darthbanana13/datastorageapi/internal/service/chat"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type saveDataUri struct {
	CustomerId uint `uri:"customerId" binding:"required"`
	DialogId   uint `uri:"dialogId" binding:"required"`
}

type saveDataPost struct {
	Text     string `json:"text" binding:"required"`
	Language string `json:"language" binding:"required"`
}

func SaveData(c *gin.Context) {
	var uriData saveDataUri
	var postData saveDataPost

	if err := c.Bind(&postData); err != nil {
		bindError(c, postData, "post", err)
		return
	}

	if err := c.BindUri(&uriData); err != nil {
		bindError(c, uriData, "uri", err)
		return
	}

	if err := chat.SaveMsg(
		uriData.CustomerId,
		uriData.DialogId,
		postData.Text,
		postData.Language,
	); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Sorry, something went wrong.\nA team of highly trained monkeys has been dispatched to deal with the situation"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"error": "none"},
	)
}

func bindError(c *gin.Context, jsonData interface{}, dataSource string, err error) {
	strVals, malformedErr := json.MarshalIndent(jsonData, "", "  ")

	if malformedErr != nil {
		log.Errorf("Could not unmarshal %s data from the request\n%s", dataSource, malformedErr)
	} else {
		log.Errorf("Could not bind %s value(s) to the request structure: %s\n%s", dataSource, strVals, err)
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid values where sent"})
}

type consentUri struct {
	DialogId uint `uri:"dialogId" binding:"required"`
}

//If it changes in the future, rerfer to https://github.com/gin-gonic/gin/issues/814 why binding required does not work with bool
type consentPost struct {
	Consent *bool `json:"consent" binding:"required"`
}

func Consent(c *gin.Context) {
	var uriData consentUri
	var postData consentPost

	if err := c.Bind(&postData); err != nil {
		bindError(c, postData, "post", err)
		return
	}

	if err := c.BindUri(&uriData); err != nil {
		bindError(c, uriData, "uri", err)
		return
	}
	chat.Consent(uriData.DialogId, *postData.Consent)
}

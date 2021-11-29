package chat

import (
	"encoding/json"
	"net/http"

	chatService "github.com/darthbanana13/datastorageapi/internal/service/chat"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type saveDataUri struct {
	CustomerId uint `uri:"customerId" binding:"required"`
	DialogId   uint `uri:"dialogId" binding:"required"`
}

//TODO: Should probably write an own validator to handle the case insensitive language
type saveDataPost struct {
	Text     string `json:"text" binding:"required"`
	Language string `json:"language" binding:"required,oneof=en fr it EN FR IT En Fr It eN fR iT"`
}

//TODO: Prevent duplicate customerId & dialogId combination from being inserted
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

	if err := chatService.SaveMsg(
		uriData.CustomerId,
		uriData.DialogId,
		postData.Text,
		postData.Language,
	); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Sorry, it's not you, it's me. Let me try again!"},
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

//If it changes in the future, rerfer to https://github.com/gin-gonic/gin/issues/814 why binding required does not work with bool but with *bool
type consentPost struct {
	Consent *bool `json:"consent" binding:"required"`
}

//TODO: Not sure if an error should be returned if the dialogId does not exist
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

	if err := chatService.Consent(uriData.DialogId, *postData.Consent); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Oh no! My spaghetti code is not working properly. I'll be back soon!"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"error": "none"},
	)
}

//TODO: Language CustomerId & DialogId could use custom validators
type viewUri struct {
	Language       string `form:"language"`
	CustomerId     uint   `form:"customerId"`
	DialogId       uint   `form:"dialogId"`
	EntriesPerPage uint   `form:"entriesPerPage,default=50" binding:"max=200"`
	Page           uint   `form:"page,default=1"`
}

func View(c *gin.Context) {
	var viewData viewUri
	if err := c.BindQuery(&viewData); err != nil {
		bindError(c, viewData, "uri", err)
		return
	}
	msgs, err := chatService.FilterAndSortByNewFirst(
		map[string]interface{}{
			"language":   viewData.Language,
			"customerId": viewData.CustomerId,
			"dialogId":   viewData.DialogId,
			"consent":    true,
		},
		viewData.Page,
		viewData.EntriesPerPage,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Oh no! Cody is sad! Internally sad!"},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		msgs,
	)
}

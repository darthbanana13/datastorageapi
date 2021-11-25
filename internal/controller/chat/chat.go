package chat

import (
	"encoding/json"
	"net/http"

	//TODO: Ideally should depend on service, not the model
	chat "github.com/darthbanana13/datastorageapi/internal/model/chat"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type saveDataUri struct {
  CustomerId	uint `uri:"customerId" binding:"required"`
  DialogId	uint `uri:"dialogId" binding:"required"`
}

type saveDataPost struct {
  Text		string `json:"text" binding:"required"`
  Language	string `json:"language" binding:"required"`
}

func SaveData(c *gin.Context) {
  var uriData saveDataUri
  var postData saveDataPost

  if err := c.Bind(&postData); err != nil {
    saveDataBindError(c, postData, "post", err)
    return
  }

  if err := c.BindUri(&uriData); err != nil {
    saveDataBindError(c, uriData, "uri", err)
    return
  }

  msg := chat.NewChatMsg(
    uriData.CustomerId,
    uriData.DialogId,
    postData.Text,
    postData.Language,
  )
  msg.Insert()

  c.JSON(
    http.StatusOK,
    gin.H{"error": "none"},
  )
}

func saveDataBindError(c *gin.Context, jsonData interface{}, dataSource string, err error) {
  strPostVals, malformedErr := json.MarshalIndent(jsonData, "", "  ")

  if malformedErr != nil {
    log.Errorf("Could not unmarshal %s data from the request\n%s", dataSource ,malformedErr)
  } else {
    log.Errorf("Could not bind a %s value(s) to the request structure: %s\n%s", dataSource, strPostVals, err)
  }
  c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid values where sent"})
}

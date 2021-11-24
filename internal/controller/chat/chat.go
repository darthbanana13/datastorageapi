package chat

import (
	"fmt"
	"strconv"
	"net/http"

	//TODO: Ideally should depend on service, not the model
	chat "github.com/darthbanana13/datastorageapi/internal/model/chat"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SaveData(c *gin.Context) {
  customerId, err := strconv.Atoi(c.Param("customerId"))
  if err != nil {
    errMsg := fmt.Sprintf("CustomerId %s is not an int value", c.Param("customerId"))
    log.Error(errMsg)
    c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
  }
  dialogId, err := strconv.Atoi(c.Param("dialogId"))
  if err != nil {
    errMsg := fmt.Sprintf("DialogId %s is not an int value", c.Param("dialogId"))
    log.Error(errMsg)
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  }

  msg := chat.NewEmptyChatMsg(customerId, dialogId)
  //TODO: This does not prevent the user from overriding the parameter values
  if err := c.ShouldBindJSON(&msg); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  }
  if err != nil {
    log.Errorf("Could not marshall data for the request")
  }
  msg.Insert()
  c.JSON(
    http.StatusOK,
    gin.H{
      "message": "received info",
      "text": msg.Text,
      "language": msg.Language,
      "customerId": c.Param("customerId"),
      "dialogId": c.Param("dialogId"),
      "error": "none",
    },
  )
}

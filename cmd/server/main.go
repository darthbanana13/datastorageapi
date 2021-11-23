package main

import (
	"os"
	"fmt"
	"net/http"

	"github.com/darthbanana13/datastorageapi/pkg/localpath"
	initLog "github.com/darthbanana13/datastorageapi/pkg/initlogrusfromtext"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type TempSaveData struct {
  Text      string `json:"text"`
  Language  string `json:"language"`
}

func init() {
  appDir, err := localpath.Get()
  if err != nil {
    log.Fatal("Could not load current running directory", err)
  }
  err = godotenv.Load(fmt.Sprintf("%s/.env", appDir))
  if err != nil {
    initLog.Init(initLog.FormatText, initLog.OutputStderr, initLog.LevelError)
    log.Fatal("Error loading .env file", err)
  }

  err = initLog.Init(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_OUTPUT"), os.Getenv("LOG_FORMAT"))
  if err != nil {
    log.Fatal(err.Error())
  }
}

func tempSaveData(c *gin.Context) {
  var chatMessage TempSaveData
  if err := c.ShouldBindJSON(&chatMessage); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  }
  c.JSON(
    http.StatusOK,
    gin.H{
      "message": "received info",
      "text": chatMessage.Text,
      "language": chatMessage.Language,
      "customerId": c.Param("customerId"),
      "dialogId": c.Param("dialogId"),
      "error": "none",
    },
  )
}

func main() {
  router := gin.Default()

  router.POST("/data/:customerId/:dialogId", tempSaveData)

  //TODO: Handle default value
  router.Run(os.Getenv("SERVER_ADDRESS"))
}

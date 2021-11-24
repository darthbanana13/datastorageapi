package main

import (
	"os"
	"fmt"
	"net/http"
	"strconv"

	"github.com/darthbanana13/datastorageapi/pkg/localpath"
	chat "github.com/darthbanana13/datastorageapi/internal/model/chat"
	arangoInit "github.com/darthbanana13/datastorageapi/pkg/initArangoDb"
	initLog "github.com/darthbanana13/datastorageapi/pkg/initlogrusfromtext"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	driver "github.com/arangodb/go-driver"
)

var db driver.Database

func init() {
  appDir, err := localpath.Get()
  if err != nil {
    log.Fatalf("Could not load current running directory: %v", err)
  }
  err = godotenv.Load(fmt.Sprintf("%s/.env", appDir))
  if err != nil {
    initLog.Init(initLog.FormatText, initLog.OutputStderr, initLog.LevelError)
    log.Fatalf("Error loading .env file: %v", err)
  }

  err = initLog.Init(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_OUTPUT"), os.Getenv("LOG_FORMAT"))
  if err != nil {
    log.Fatalf("Error initializing log settings: %v", err)
  }

  db, err = arangoInit.InitDbWith(
    os.Getenv("ARANGODB_PROTOCOL") + os.Getenv("ARANGODB_HOST") + ":" + os.Getenv("ARANGODB_PORT"),
    os.Getenv("ARANGODB_USER"),
    os.Getenv("ARANGODB_PASSWORD"),
    os.Getenv("ARANGODB_NAME"),
    []string{"chat"},
  )

  if err != nil {
    log.Fatalf("Error initializing db: %v", err)
  }
}

func insertArangoDb(jsonData string) {
  db.Query(nil, "INSERT" + jsonData + "INTO chat", nil)
}

func tempSaveData(c *gin.Context) {
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
  msg.Insert(db)
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

func main() {
  router := gin.Default()

  router.POST("/data/:customerId/:dialogId", tempSaveData)

  //TODO: Handle default value
  router.Run(os.Getenv("SERVER_ADDRESS"))
}

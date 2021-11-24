package main

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"strconv"

	"github.com/darthbanana13/datastorageapi/pkg/localpath"
	arangoInit "github.com/darthbanana13/datastorageapi/pkg/initArangoDb"
	initLog "github.com/darthbanana13/datastorageapi/pkg/initlogrusfromtext"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	driver "github.com/arangodb/go-driver"
)

type TempSaveData struct {
  Text		    string `json:"text"`
  Language	    string `json:"language"`
  CustomerId	    int
  DialogId	    int
  Consent	    bool
  NanoTimestamp     int64
}

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
}

func insertArangoDb(jsonData string) {
  db.Query(nil, "INSERT" + jsonData + "INTO chat", nil)
}

//TODO: Return error if URL parameters are not int
func tempSaveData(c *gin.Context) {
  var chatMessage TempSaveData
  if err := c.ShouldBindJSON(&chatMessage); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  }
  chatMessage.CustomerId, _ = strconv.Atoi(c.Param("customerId"))
  chatMessage.DialogId, _ = strconv.Atoi(c.Param("dialogId"))
  chatMessage.Consent = false
  chatMessage.NanoTimestamp = time.Now().UnixNano()
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

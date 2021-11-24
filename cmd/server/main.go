package main

import (
  "os"

  "github.com/darthbanana13/datastorageapi/internal/initApp"
  "github.com/darthbanana13/datastorageapi/internal/controller/chat"

  "github.com/gin-gonic/gin"
)

func init() {
  initApp.InitAll()
}

func main() {
  router := gin.Default()

  router.POST("/data/:customerId/:dialogId", chat.SaveData)

  //TODO: Handle default value
  router.Run(os.Getenv("SERVER_ADDRESS"))
}

package main

import (
    "os"
    "fmt"

    "github.com/darthbanana13/datastorageapi/pkg/localpath"
    initLog "github.com/darthbanana13/datastorageapi/pkg/initlogrusfromtext"

    log "github.com/sirupsen/logrus"
    "github.com/joho/godotenv"
)

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

func main() {
}

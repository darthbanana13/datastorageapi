package initApp

import (
  "os"

  "github.com/darthbanana13/datastorageapi/pkg/localpath"
  arangoInit "github.com/darthbanana13/datastorageapi/pkg/initArangoDb"

  "github.com/joho/godotenv"
  log "github.com/sirupsen/logrus"
  driver "github.com/arangodb/go-driver"
  container "github.com/golobby/container/v3"
  initLogrus "github.com/darthbanana13/datastorageapi/pkg/initlogrusfromtext"
)

func initEnv() {
  appDir, err := localpath.Get()
  if err != nil {
    log.Fatalf("Could not load current running directory: %v", err)
  }
  err = godotenv.Load(appDir + "/.env")
  if err != nil {
    initLogrus.Init(initLogrus.FormatText, initLogrus.OutputStderr, initLogrus.LevelError)
    log.Fatalf("Error loading .env file: %v", err)
  }
}

func initLog() {
  err := initLogrus.Init(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_OUTPUT"), os.Getenv("LOG_FORMAT"))
  if err != nil {
    log.Fatalf("Error initializing log settings: %v", err)
  }
}

func initDb() {
  db, err := arangoInit.InitDbWith(
    os.Getenv("ARANGODB_PROTOCOL") + os.Getenv("ARANGODB_HOST") + ":" + os.Getenv("ARANGODB_PORT"),
    os.Getenv("ARANGODB_USER"),
    os.Getenv("ARANGODB_PASSWORD"),
    os.Getenv("ARANGODB_NAME"),
    []string{"chat"},
  )

  if err != nil {
    log.Fatalf("Error initializing db: %v", err)
  }
  err = container.Singleton(func() *driver.Database {
      return &db
  })
}

func InitAll() {
  initEnv()
  initLog()
  initDb()
}

package main

import (
    "os"
    "fmt"
    "strings"

    "github.com/darthbanana13/datastorageapi/pkg/localpath"

    log "github.com/sirupsen/logrus"
    "github.com/joho/godotenv"
)

const format_text = "text"
const format_json = "json"

const output_stdout = "stdout"
const output_stderr = "stderr"

const level_trace = "trace"
const level_debug = "debug"
const level_info = "info"
const level_warn = "warn"
const level_error = "error"
const level_fatal = "fatal"
const level_panic = "panic"

func initLogFormat(format string) error {
  lowerFormat := strings.ToLower(format)
  switch lowerFormat {
  case format_text:
    log.SetFormatter(&log.TextFormatter{})
  case format_json:
    log.SetFormatter(&log.JSONFormatter{})
  default:
    log.SetFormatter(&log.TextFormatter{})
    return fmt.Errorf("Unknown log format %s configured for LOG_FORMAT", format);
  }
  return nil;
}

func initLogOutput(output string) error {
  lowerOutput := strings.ToLower(output)
  switch lowerOutput {
  case output_stdout:
    log.SetOutput(os.Stdout)
  case output_stderr:
    log.SetOutput(os.Stderr)
  default:
    log.SetOutput(os.Stderr)
    return fmt.Errorf("Unknown log output %s configured for LOG_OUTPUT", output)
  }
  return nil
}

func initLogLevel(level string) error {
  lowerLevel := strings.ToLower(level)
  switch lowerLevel {
  case level_trace:
    log.SetLevel(log.TraceLevel)
  case level_debug:
    log.SetLevel(log.DebugLevel)
  case level_info:
    log.SetLevel(log.InfoLevel)
  case level_warn:
    log.SetLevel(log.WarnLevel)
  case level_error:
    log.SetLevel(log.ErrorLevel)
  case level_fatal:
    log.SetLevel(log.FatalLevel)
  case level_panic:
    log.SetLevel(log.PanicLevel)
  default:
    log.SetOutput(os.Stderr)
    return fmt.Errorf("Unknown log level %s configured for LOG_LEVEL", level)
  }
  return nil
}

// TODO: Should return all errors
func initLogging(level, output, format string) error {
  err := initLogFormat(format)
  if err != nil {
    return err
  }
  initLogOutput(output)
  if err != nil {
    return err
  }
  return initLogLevel(level)
}

func init() {
  appDir, _ := localpath.Get()
  err := godotenv.Load(fmt.Sprintf("%s/.env", appDir))
  if err != nil {
    initLogging(format_text, output_stderr, level_error)
    log.Fatal("Error loading .env file")
  }

  err = initLogging(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_OUTPUT"), os.Getenv("LOG_FORMAT"))
  if err != nil {
    log.Fatal(err.Error())
  }
}

func main() {
}

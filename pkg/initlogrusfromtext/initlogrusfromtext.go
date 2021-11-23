package initlogrusfromtext

import (
    "os"
    "fmt"
    "strings"

    "github.com/sirupsen/logrus"
)

const FormatText = "text"
const FormatJson = "json"

const OutputStdout = "stdout"
const OutputStderr = "stderr"

const LevelTrace = "trace"
const LevelDebug = "debug"
const LevelInfo = "info"
const LevelWarn = "warn"
const LevelError = "error"
const LevelFatal = "fatal"
const LevelPanic = "panic"

func InitLogFormat(format string) error {
  lowerFormat := strings.ToLower(format)
  switch lowerFormat {
  case FormatText:
    logrus.SetFormatter(&logrus.TextFormatter{})
  case FormatJson:
    logrus.SetFormatter(&logrus.JSONFormatter{})
  default:
    logrus.SetFormatter(&logrus.TextFormatter{})
    return fmt.Errorf("Unknown log format %s configured for LOG_FORMAT", format);
  }
  return nil;
}

func InitLogOutput(output string) error {
  lowerOutput := strings.ToLower(output)
  switch lowerOutput {
  case OutputStdout:
    logrus.SetOutput(os.Stdout)
  case OutputStderr:
    logrus.SetOutput(os.Stderr)
  default:
    logrus.SetOutput(os.Stderr)
    return fmt.Errorf("Unknown log output %s configured for LOG_OUTPUT", output)
  }
  return nil
}

func InitLogLevel(level string) error {
  lowerLevel := strings.ToLower(level)
  switch lowerLevel {
  case LevelTrace:
    logrus.SetLevel(logrus.TraceLevel)
  case LevelDebug:
    logrus.SetLevel(logrus.DebugLevel)
  case LevelInfo:
    logrus.SetLevel(logrus.InfoLevel)
  case LevelWarn:
    logrus.SetLevel(logrus.WarnLevel)
  case LevelError:
    logrus.SetLevel(logrus.ErrorLevel)
  case LevelFatal:
    logrus.SetLevel(logrus.FatalLevel)
  case LevelPanic:
    logrus.SetLevel(logrus.PanicLevel)
  default:
    logrus.SetOutput(os.Stderr)
    return fmt.Errorf("Unknown log level %s configured for LOG_LEVEL", level)
  }
  return nil
}

// TODO: Should return all errors
func Init(level, output, format string) error {
  err := InitLogFormat(format)
  if err != nil {
    return err
  }
  InitLogOutput(output)
  if err != nil {
    return err
  }
  return InitLogLevel(level)
}

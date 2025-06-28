package logger

import (
	"log"
	"os"

	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/utils"
)

var logLevel = utils.GetEnv("LOG_LEVEL", "DEBUG")

type Log interface {
	CustomDebug(title string, message interface{}, options ...*Options)
	CustomError(title string, message interface{}, options ...*Options)
	CustomLog(title string, message interface{}, options ...*Options)
	CustomPanic(title string, message interface{}, options ...Options)
	Debug(message interface{}, options ...*Options)
	Error(message interface{}, options ...*Options)
	Log(message interface{}, options ...*Options)
	Panic(message interface{}, options ...Options)
}

type logImpl struct {
	stdout *log.Logger
	stderr *log.Logger
}

type Options struct {
	IsPrintStack bool
	IsExit       bool
	ExitCode     int
}

func New(prefix string) Log {
	return &logImpl{
		stdout: log.New(os.Stdout, "[LOG]["+prefix+"]", log.Ldate|log.Ltime),
		stderr: log.New(os.Stderr, "[ERROR]["+prefix+"]", log.Ldate|log.Ltime),
	}
}

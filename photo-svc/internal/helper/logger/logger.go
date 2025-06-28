package logger

import (
	"log"
	"os"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/utils"
)

var logLevel = utils.GetEnv("LOG_LEVEL")

type Log struct {
	stdout *log.Logger
	stderr *log.Logger
}

type Options struct {
	IsPrintStack bool
	IsExit       bool
	ExitCode     int
}

func New(prefix string) *Log {
	return &Log{
		stdout: log.New(os.Stdout, "[LOG]["+prefix+"]", log.Ldate|log.Ltime),
		stderr: log.New(os.Stderr, "[ERROR]["+prefix+"]", log.Ldate|log.Ltime),
	}
}

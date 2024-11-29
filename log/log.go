package log

import (
	"os"
	"strings"
	"syscall"

	"github.com/op/go-logging"
)

const (
	END = "\033[0m"

	BLACK  = "\033[30m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PURPLE = "\033[35m"
	CYAN   = "\033[36m"
	GREY   = "\033[90m"

	HIGH_RED    = "\033[91m"
	HIGH_GREEN  = "\033[92m"
	HIGH_YELLOW = "\033[93m"
	HIGH_BLUE   = "\033[94m"
	HIGH_PURPLE = "\033[95m"
	HIGH_CYAN   = "\033[96m"
)

var Logger = logging.MustGetLogger("logger")

var fileFormat = logging.MustStringFormatter(
	`%{time:15:04:05.000} %{level:.8s} ▶ %{message}`,
)
var backendFormat = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{level:.8s} ▶%{color:reset} %{message}`,
)

var (
	backends []logging.Backend
	logFile  *os.File
)

func init() {
	backends = make([]logging.Backend, 0, 1)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, backendFormat)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.INFO, "")

	backends = append(backends, backendLeveled)
	logging.SetBackend(backends...)
}

func SetLogFile(file string) {
	var err error

	logFile, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	err = syscall.Dup2(int(logFile.Fd()), int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	filebackend := logging.NewLogBackend(os.Stdout, "", 0)
	filebackendFormatter := logging.NewBackendFormatter(filebackend, fileFormat)
	filebackendLeveled := logging.AddModuleLevel(filebackendFormatter)
	filebackendLeveled.SetLevel(logging.INFO, "")

	backends = append(backends, filebackendLeveled)
	logging.SetBackend(backends...)
}

func CloseLogFile() {
	logFile.Close()
}

func SetLogLevel(newLevel string) {
	var level logging.Level
	lev := strings.ToLower(newLevel)
	switch lev {
	case "debug":
		level = logging.DEBUG
	default:
		fallthrough
	case "info":
		level = logging.INFO
	case "notice":
		level = logging.NOTICE
	case "warning":
		level = logging.WARNING
	case "error":
		level = logging.ERROR
	case "critical":
		level = logging.CRITICAL
	}
	logging.SetLevel(level, "")
}

var Debug = Logger.Debug
var Debugf = Logger.Debugf
var Info = Logger.Info
var Infof = Logger.Infof
var Notice = Logger.Notice
var Noticef = Logger.Noticef
var Warning = Logger.Warning
var Warningf = Logger.Warningf
var Error = Logger.Error
var Errorf = Logger.Errorf
var Critical = Logger.Critical
var Criticalf = Logger.Criticalf
var Fatal = Logger.Fatal
var Fatalf = Logger.Fatalf
var Panic = Logger.Panic
var Panicf = Logger.Panicf

package logs

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type LogLevel int

const (
	Trace = iota
	Info
	Warning
	Error
	Fatal
)

type ErrorLogger interface {
	PrintError(level LogLevel, err error)
	TryPrintError(level LogLevel, err error)
	TryTrace(err error)
	TryInfo(err error)
	TryWarning(err error)
	TryError(err error)
	TryFatal(err error)
}

type Logger interface {
	ErrorLogger
	Println(level LogLevel, message string)
	Printlnf(level LogLevel, format string, a ...interface{})
	Trace(message string)
	Tracef(format string, a ...interface{})
	Info(message string)
	Infof(format string, a ...interface{})
	Warning(message string)
	Warningf(format string, a ...interface{})
	Error(message string)
	Errorf(format string, a ...interface{})
	Fatal(message string)
	Fatalf(format string, a ...interface{})
	SetConsoleLevel(level LogLevel)
	SetFileLogLevel(level LogLevel)
	GetErrorLogger() *log.Logger
	SetDir(string)
}

type logger struct {
	loggers             map[LogLevel]*log.Logger
	activeConsoleLogger map[LogLevel]bool
	activeFileLogger    map[LogLevel]bool
	errorLogger         *log.Logger
	logsDir             string
}

func NewLogger() Logger {
	logger := &logger{loggers: make(map[LogLevel]*log.Logger), activeConsoleLogger: setLevel(Trace), activeFileLogger: setLevel(Trace)}
	createLog := func(level LogLevel, writer io.Writer) *log.Logger {
		levels := [...]string{
			"TRACE: ",
			"INFO: ",
			"WARNING: ",
			"ERROR: ",
			"FATAL: "}
		return log.New(writer,
			levels[level],
			log.Ldate|log.Ltime)
	}
	registerLogger := func(consoleWriter io.Writer, levels ...LogLevel) {
		for _, level := range levels {
			logger.loggers[level] = createLog(level, consoleWriter)
		}
	}

	registerLogger(os.Stdout, Trace, Info, Warning)
	registerLogger(os.Stderr, Error, Fatal)
	logger.errorLogger = logger.loggers[Error]

	return logger
}

func (logger *logger) Printlnf(level LogLevel, format string, a ...interface{}) {
	logger.Println(level, fmt.Sprintf(format, a...))
}

func (logger *logger) Println(level LogLevel, message string) {
	var writers []io.Writer
	if logger.activeConsoleLogger[level] {
		if level < Error {
			writers = append(writers, os.Stdout)
		} else {
			writers = append(writers, os.Stderr)
		}
	}

	if len(logger.logsDir) > 0 && logger.activeFileLogger[level] {
		year, month, day := time.Now().Date()
		execFile := filepath.Base(os.Args[0])

		logFile := logger.logsDir + "/" + execFile + "-" + strconv.Itoa(year) + strconv.Itoa(int(month)) + strconv.Itoa(day) + ".log"
		_ = os.MkdirAll(filepath.Dir(logFile), 0666)
		osFile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println("impossible to log in file:", err)
		}
		defer func() {
			_ = osFile.Close()
		}()
		writers = append(writers, osFile)
	}

	if writers != nil {
		multiWriter := io.MultiWriter(writers...)
		logger.loggers[level].SetOutput(multiWriter)

		if level == Fatal {
			logger.loggers[level].Fatalln(message)
		} else {
			logger.loggers[level].Println(message)
		}
	}
}

func (logger *logger) SetDir(directory string) {
	logger.logsDir = directory
}

func (logger *logger) SetConsoleLevel(level LogLevel) {
	logger.activeConsoleLogger = setLevel(level)
}

func (logger *logger) SetFileLogLevel(level LogLevel) {
	logger.activeFileLogger = setLevel(level)
}

func (logger *logger) GetErrorLogger() *log.Logger {
	return logger.errorLogger
}

func (logger *logger) PrintError(level LogLevel, err error) {
	errorschecker.CheckNilParameter(map[string]interface{}{"err": err})
	logger.Println(level, err.Error())
}

func (logger *logger) TryPrintError(level LogLevel, err error) {
	if err != nil {
		logger.PrintError(level, err)
	}
}

func (logger *logger) Trace(message string) {
	logger.Println(Trace, message)
}

func (logger *logger) Tracef(format string, a ...interface{}) {
	logger.Printlnf(Trace, format, a...)
}

func (logger *logger) TryTrace(err error) {
	logger.TryPrintError(Trace, err)
}

func (logger *logger) Info(message string) {
	logger.Println(Info, message)
}
func (logger *logger) Infof(format string, a ...interface{}) {
	logger.Printlnf(Info, format, a...)
}

func (logger *logger) TryInfo(err error) {
	logger.TryPrintError(Info, err)
}

func (logger *logger) Warning(message string) {
	logger.Println(Warning, message)
}

func (logger *logger) Warningf(format string, a ...interface{}) {
	logger.Printlnf(Warning, format, a...)
}

func (logger *logger) TryWarning(err error) {
	logger.PrintError(Warning, err)
}

func (logger *logger) Error(message string) {
	logger.Println(Error, message)
}

func (logger *logger) Errorf(format string, a ...interface{}) {
	logger.Printlnf(Error, format, a...)
}

func (logger *logger) TryError(err error) {
	logger.TryPrintError(Error, err)
}

func (logger *logger) Fatal(message string) {
	logger.Println(Fatal, message)
}

func (logger *logger) Fatalf(format string, a ...interface{}) {
	logger.Printlnf(Fatal, format, a...)
}

func (logger *logger) TryFatal(err error) {
	logger.PrintError(Fatal, err)
}

func setLevel(level LogLevel) map[LogLevel]bool {
	loggersActives := map[LogLevel]bool{Trace: true, Info: true, Warning: true, Error: true, Fatal: true}
	if level > Trace {
		start := level - 1
		for i := start; i > -1; i-- {
			loggersActives[i] = false
		}
	}
	return loggersActives
}

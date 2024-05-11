package logwrapper

import (
	"log"
	"os"
	"strings"
	"time"
)

type FileLogger struct {
	l              *log.Logger
	EnableRotate   bool
	DirLogger      string
	PrefixFileName string
	current        time.Time
	file           *os.File
	innerLogger    ILoggerWrapper
}

func (logger *FileLogger) Log(msgs ...ILogMsg) {
	if logger.EnableRotate {
		nYear, nMonth, nDay := time.Now().Date()
		cYear, cMonth, cDay := logger.current.Date()
		if nYear != cYear || nMonth != cMonth || nDay != cDay {
			logger.file.Close()
			newPath := strings.Join(
				[]string{logger.DirLogger, logger.PrefixFileName + "-" + time.Now().Format(time.DateOnly)},
				"/",
			) + ".log"
			newWriter, err := os.OpenFile(newPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Panic(err)
			}
			logger.file = newWriter
			logger.l.SetOutput(newWriter)
		}
	}
	for _, v := range msgs {
		logger.l.Println(v.String())
	}
	if logger.innerLogger != nil {
		logger.innerLogger.Log(msgs...)
	}
}

func (logger *FileLogger) AddLogger(l ILoggerWrapper) {
	logger.innerLogger = l
}

func NewFileLoggerInstance(enableRotate bool, dir string, prefix string) *FileLogger {
	defaultPath := strings.Join(
		[]string{dir, prefix + "-" + time.Now().Format(time.DateOnly)},
		"/",
	) + ".log"
	defaultWriter, err := os.OpenFile(defaultPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	var fileLogger = FileLogger{
		l:              log.New(defaultWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		current:        time.Now(),
		EnableRotate:   enableRotate,
		DirLogger:      dir,
		PrefixFileName: prefix,
		file:           defaultWriter,
		innerLogger:    nil,
	}
	return &fileLogger
}

type StandarLogger struct {
	l           *log.Logger
	innerLogger ILoggerWrapper
}

func (logger *StandarLogger) Log(msgs ...ILogMsg) {
	for _, v := range msgs {
		logger.l.Println(v.String())
	}
	if logger.innerLogger != nil {
		logger.innerLogger.Log(msgs...)
	}
}

func (logger *StandarLogger) AddLogger(l ILoggerWrapper) {
	logger.innerLogger = l
}

func NewStanderLoggerInstance() *StandarLogger {
	return &StandarLogger{
		l:           log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
		innerLogger: nil,
	}
}

type LoggerWrapperBuilder struct {
	logger ILoggerWrapper
}

func (builder *LoggerWrapperBuilder) AddLogger(l ILoggerWrapper) *LoggerWrapperBuilder {
	if builder.logger != nil {
		l.AddLogger(builder.logger)
		return &LoggerWrapperBuilder{logger: l}
	}
	builder.logger = l
	return builder
}

func (builder *LoggerWrapperBuilder) Build() ILoggerWrapper {
	return builder.logger
}

func NewLoggerWrapperBuilder() *LoggerWrapperBuilder {
	return &LoggerWrapperBuilder{}
}

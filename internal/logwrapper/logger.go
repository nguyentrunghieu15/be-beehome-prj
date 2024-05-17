package logwrapper

import (
	"io"
	"log"

	"github.com/natefinch/lumberjack"
)

type ConfigRollbackWriter struct {
	MaxAge, MaxSize, MaxBackups int
	Compress                    bool
}

func NewRollbackWriterFile(logFilename string, config ConfigRollbackWriter) io.Writer {
	writer := &lumberjack.Logger{
		Filename:   logFilename,
		MaxAge:     config.MaxAge,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
	}
	return writer
}

type ILoggerWrapper interface {
	Infor(string)
	Error(string)
	Warn(string)
	Debug(string)
	Printf(string, ...interface{})
	SetWriter(io.Writer)
}

// LoggerWrapper struct implements the ILoggerWrapper interface
type LoggerWrapper struct {
	logger *log.Logger
}

// NewLoggerWrapper creates a new LoggerWrapper instance with an optional prefix
func (*LoggerWrapper) Init() interface{} {
	return &LoggerWrapper{
		logger: log.New(
			log.Writer(),
			"",
			log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix),
	}
}

// Info logs a message with INFO level
func (l *LoggerWrapper) Info(msg string) {
	l.log(" [info] " + msg)
}

// Error logs a message with ERROR level
func (l *LoggerWrapper) Error(msg string) {
	l.log(" [error] " + msg)
}

// Warn logs a message with WARN level
func (l *LoggerWrapper) Warn(msg string) {
	l.log(" [warn] " + msg)
}

// Debug logs a message with DEBUG level
func (l *LoggerWrapper) Debug(msg string) {
	l.log(" [debug] " + msg)
}

// Printf logs a formatted message with any level based on the log package flags
func (l *LoggerWrapper) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *LoggerWrapper) log(msg string) {
	l.logger.Println(msg)
}

func (l *LoggerWrapper) SetWriter(writer io.Writer) {
	l.logger.SetOutput(writer)
}

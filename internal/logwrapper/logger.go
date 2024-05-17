package logwrapper

import (
	"io"
	"log"

	"github.com/natefinch/lumberjack"
)

func NewRollbackWriterFile(logFilename string, maxAge, maxSize, maxBackups int) io.Writer {
	writer := &lumberjack.Logger{
		Filename:   logFilename,
		MaxAge:     maxAge,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		Compress:   true,
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
func NewLoggerWrapper() *LoggerWrapper {
	return &LoggerWrapper{
		logger: log.New(
			log.Writer(),
			"",
			log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix),
	}
}

// Info logs a message with INFO level
func (l *LoggerWrapper) Info(msg string) {
	l.log("| Info | " + msg)
}

// Error logs a message with ERROR level
func (l *LoggerWrapper) Error(msg string) {
	l.log("| Error | " + msg)
}

// Warn logs a message with WARN level
func (l *LoggerWrapper) Warn(msg string) {
	l.log("| Warn | " + msg)
}

// Debug logs a message with DEBUG level
func (l *LoggerWrapper) Debug(msg string) {
	l.log("| Debug | " + msg)
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

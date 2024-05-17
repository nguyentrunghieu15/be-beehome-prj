package logwrapper

import (
	"testing"
)

func TestLoggerWrapper(t *testing.T) {
	// Test cases for different log levels
	testCases := []struct {
		level    string
		message  string
		expected string
	}{
		{"Info", "Application started", "| Info | Application started\n"},
		{"Error", "Error occurred", "| Error | Error occurred\n"},
		{"Warn", "Potential issue detected", "| Warn | Potential issue detected\n"},
		{"Debug", "Detailed debug information", "| Debug | Detailed debug information\n"},
	}

	// Create a logger with a prefix
	logger := NewLoggerWrapper()

	// Call logging methods for each test case
	for _, tc := range testCases {
		switch tc.level {
		case "Info":
			logger.Info(tc.message)
		case "Error":
			logger.Error(tc.message)
		case "Warn":
			logger.Warn(tc.message)
		case "Debug":
			logger.Debug(tc.message)
		}
	}
}

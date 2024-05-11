package logwrapper

import (
	"testing"
)

func TestFileLogger_Log(t *testing.T) {

	tests := []struct {
		name   string
		logger *FileLogger
	}{
		{
			name:   "Standar",
			logger: NewFileLoggerInstance(true, ".", "test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.Log(NewInforMsg("Test"))
		})
	}
}

func TestStandarLogger_Log(t *testing.T) {
	tests := []struct {
		name   string
		logger *StandarLogger
	}{
		{
			name:   "Standar",
			logger: NewStanderLoggerInstance(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.Log(NewInforMsg("Test"))
		})
	}
}

func TestLoggerWrapperBuilder_Build(t *testing.T) {

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "Must_pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewLoggerWrapperBuilder().AddLogger(NewFileLoggerInstance(false, ".", "test")).
				AddLogger(NewStanderLoggerInstance()).Build().Log(NewInforMsg("Test"))
		})
	}
}

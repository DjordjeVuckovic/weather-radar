package logger

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestInitSlog(t *testing.T) {
	tests := []struct {
		name            string
		config          Config
		expectedHandler Handler
	}{
		{"DebugLevel with Text", Config{Level: DebugLevel, Handler: Text}, Text},
		{"InfoLevel with JSON", Config{Level: InfoLevel, Handler: Json}, Json},
		{"WarnLevel with Text", Config{Level: WarnLevel, Handler: Text}, Text},
		{"ErrorLevel with JSON", Config{Level: ErrorLevel, Handler: Json}, Json},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer, _ := os.Pipe()
			originalStdout := os.Stdout
			defer func() {
				os.Stdout = originalStdout
			}()
			os.Stdout = writer

			var buf bytes.Buffer
			done := make(chan struct{})
			go func() {
				_, _ = io.Copy(&buf, reader)
				close(done)
			}()

			InitSlog(tt.config)

			slog.Info("info message")
			slog.Warn("warn message")
			slog.Error("error message")

			_ = writer.Close()
			<-done

			output := buf.String()
			if tt.config.Handler == Json && !isJSONFormat(output) {
				t.Errorf("Expected JSON format, got text format for %s", tt.name)
			}
			if tt.config.Handler == Text && isJSONFormat(output) {
				t.Errorf("Expected text format, got JSON format for %s", tt.name)
			}
		})
	}
}

// Helper function to check if the output is in JSON format
func isJSONFormat(output string) bool {
	return len(output) > 0 && output[0] == '{'
}

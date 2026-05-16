package main

import (
	"log/slog"
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want slog.Level
	}{
		{name: "debug", raw: "debug", want: slog.LevelDebug},
		{name: "info", raw: "info", want: slog.LevelInfo},
		{name: "warn", raw: "warn", want: slog.LevelWarn},
		{name: "warning alias", raw: "warning", want: slog.LevelWarn},
		{name: "error", raw: "error", want: slog.LevelError},
		{name: "trim and lower", raw: " INFO ", want: slog.LevelInfo},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseLogLevel(tt.raw)
			if err != nil {
				t.Fatalf("parseLogLevel(%q) returned error: %v", tt.raw, err)
			}

			if got != tt.want {
				t.Fatalf("parseLogLevel(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestParseLogLevelRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	if _, err := parseLogLevel("trace"); err == nil {
		t.Fatal("parseLogLevel() expected error for invalid value")
	}
}

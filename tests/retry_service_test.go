package tests

import "testing"

func TestPlaceholderRetryService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "first retry"},
		{name: "second retry"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
	}
}
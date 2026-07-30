package tests

import "testing"

func TestPlaceholderAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration placeholder in short mode")
	}
}
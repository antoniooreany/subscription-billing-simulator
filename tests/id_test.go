package tests

import (
	"strings"
	"testing"

	"github.com/antoniooreany/subscription-billing-simulator/internal/shared"
)

func TestNewIDHasPrefix(t *testing.T) {
	id := shared.NewID("cust")
	if !strings.HasPrefix(id, "cust_") {
		t.Fatalf("expected prefix cust_, got %s", id)
	}
}

func TestNewIDProducesDifferentValues(t *testing.T) {
	id1 := shared.NewID("evt")
	id2 := shared.NewID("evt")

	if id1 == id2 {
		t.Fatalf("expected different ids, got identical values: %s", id1)
	}
}

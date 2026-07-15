package reman

import "testing"

// TestRMCCEmptyPayloads verifies RMCCEmptyPayloads behavior.
func TestRMCCEmptyPayloads(t *testing.T) {
	if QueryIDPayload() != nil { t.Fatal("query id payload should be nil") }
	if PingPayload() != nil { t.Fatal("ping payload should be nil") }
}

// tracker_test.go
package main

import (
	"strings"
	"testing"
)

func TestDistractionDetection(t *testing.T) {
	mock := []string{"chrome.exe - YouTube", "code.exe"}
	distracted := false

	for _, proc := range mock {
		if strings.Contains(proc, "YouTube") || strings.Contains(proc, "Instagram") {
			distracted = true
			break
		}
	}

	if !distracted {
		t.Error("Expected distraction not detected")
	}
}

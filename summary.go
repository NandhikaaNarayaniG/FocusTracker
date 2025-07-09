// summary.go
package main

import (
	"fmt"
	"os"
	"strings"
)

func GenerateSummary() string {
	var focusCount, distractCount int
	var totalFocusedSeconds, totalDistractedSeconds int

	mu.Lock()
	defer mu.Unlock()

	if len(logLines) == 0 {
		return "No data to analyze.\n"
	}

	summary := "\n\nFocusFlow Report:\n"

	for _, line := range logLines {
		parts := strings.SplitN(line, "] ", 2)
		if len(parts) < 2 {
			continue
		}
		status := parts[1]

		if strings.Contains(status, "Focused") {
			focusCount++
			totalFocusedSeconds += 10
		} else {
			distractCount++
			totalDistractedSeconds += 10
		}
	}

	total := focusCount + distractCount
	if total == 0 {
		return "No data to analyze.\n"
	}

	focusPct := float64(focusCount) / float64(total) * 100
	distractPct := 100 - focusPct

	focusMin := totalFocusedSeconds / 60
	focusSec := totalFocusedSeconds % 60
	distractMin := totalDistractedSeconds / 60
	distractSec := totalDistractedSeconds % 60

	summary += fmt.Sprintf("Total Time Points Tracked: %d\n", total)
	summary += fmt.Sprintf(" Focused: %d (%.1f%%) - %d minutes %d seconds\n", focusCount, focusPct, focusMin, focusSec)
	summary += fmt.Sprintf("Distracted: %d (%.1f%%) - %d minutes %d seconds\n", distractCount, distractPct, distractMin, distractSec)

	return summary
}

func SaveSummaryToFile(summary string) error {
	f, err := os.OpenFile("report.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	summary = "\n\n" + summary
	_, err = f.WriteString(summary)
	return err
}

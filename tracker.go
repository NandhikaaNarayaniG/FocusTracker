// tracker.go
package main

import (
	"strings"
	"sync"
	"time"
)

func TrackDistraction(statusChannel chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		processes := GetMockProcesses()

		hasChrome := false
		hasYouTube := false
		hasInstagram := false

		for _, proc := range processes {
			if strings.Contains(proc, "Chrome") {
				hasChrome = true
			}
			if strings.Contains(proc, "YouTube") {
				hasYouTube = true
			}
			if strings.Contains(proc, "Instagram") {
				hasInstagram = true
			}
		}

		var distractions []string
		if hasChrome && hasYouTube {
			distractions = append(distractions, "YouTube")
		}
		if hasChrome && hasInstagram {
			distractions = append(distractions, "Instagram")
		}

		if len(distractions) > 0 {
			statusChannel <- "Distracted (" + strings.Join(distractions, ", ") + ")"
		} else {
			statusChannel <- "Focused"
		}

		time.Sleep(10 * time.Second)
	}
}

// helper function to avoid duplicates
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

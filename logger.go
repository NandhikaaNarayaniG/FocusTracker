package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var logLines []string
var mu sync.Mutex

func LogStatus(statusChannel <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.Create("report.txt")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer f.Close()

	for status := range statusChannel {
		timestamp := time.Now().Format("03:04:05 PM")
		log := fmt.Sprintf("[%s] %s", timestamp, status)

		// 1. Print to CLI
		fmt.Println(log)

		// 2. Save in memory for summary (if needed)
		mu.Lock()
		logLines = append(logLines, log)
		mu.Unlock()

		// 3. Write to file
		f.WriteString(log + "\n")

		// 4. Send to frontend via SSE
		logStream <- log
	}
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

var (
	statusChannel = make(chan string)
	logStream     = make(chan string) // Used to stream logs to frontend
	wg            sync.WaitGroup
)

func main() {
	// Serve static assets
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/start", startTracking)
	http.HandleFunc("/logs", streamLogs) // SSE log endpoint

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func startTracking(w http.ResponseWriter, r *http.Request) {
	go func() {
		logStream <- "🔵 Tracking started... Press Ctrl+C to stop.\n"

		wg.Add(1)
		go TrackDistraction(statusChannel, &wg)

		wg.Add(1)
		go LogStatus(statusChannel, &wg)

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan

		logStream <- "\n🟠 Stopping FocusFlow...\n"

		close(statusChannel)
		wg.Wait()

		summary := GenerateSummary()
		err := SaveSummaryToFile(summary)
		if err != nil {
			logStream <- fmt.Sprintf(" Error writing summary: %v\n", err)
		} else {
			logStream <- "Summary added to report.txt\n"
		}
	}()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func streamLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Stream logs as they arrive
	for msg := range logStream {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

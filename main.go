package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// https://docs.gitlab.com/ee/administration/audit_event_schema.html#audit-event-json-schema
type Event struct {
	ID            string      `json:"id"`
	AuthorID      int         `json:"author_id"`
	AuthorName    string      `json:"author_name"`
	Details       interface{} `json:"details"`
	IPAddress     string      `json:"ip_address"`
	EntityID      int         `json:"entity_id"`
	EntityPath    string      `json:"entity_path"`
	EntityType    string      `json:"entity_type"`
	EventType     string      `json:"event_type"`
	TargetID      int         `json:"target_id"`
	TargetType    string      `json:"target_type"`
	TargetDetails string      `json:"target_details"`
}

func main() {
	logFile, err := os.OpenFile("server.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api", apiHandler)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! you've requested %s\n", r.URL.Path)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	logData := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"payload":   string(body),
	}
	logEntry, err := json.Marshal(logData)
	if err != nil {
		http.Error(w, "Failed to serialize log data", http.StatusInternalServerError)
		return
	}
	logFile, err := os.OpenFile("server.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer logFile.Close()

	_, err = logFile.WriteString(string(logEntry) + "\n")
	if err != nil {
		http.Error(w, "Failed to write to log file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JSON data received successfully\n"))
}

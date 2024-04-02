package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Define a struct to represent the JSON data
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
	// Open a log file in append mode
	logFile, err := os.OpenFile("server.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	// Set log output to the log file
	log.SetOutput(logFile)

	// Register handler for "/json" endpoint
	http.HandleFunc("/json", jsonHandler)

	// Start the server
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Close the request body
	defer r.Body.Close()

	// Write the request body JSON to the log file
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

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JSON data received successfully\n"))
}

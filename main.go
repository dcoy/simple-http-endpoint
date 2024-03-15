package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Request struct {
	Properties struct {
		ID            string                 `json:"id"`
		AuthorID      int                    `json:"author_id"`
		AuthorName    string                 `json:"author_name"`
		Details       map[string]interface{} `json:"details"`
		IPAddress     string                 `json:"ip_address"`
		EntityID      int                    `json:"entity_id"`
		EntityPath    string                 `json:"entity_path"`
		EntityType    string                 `json:"entity_type"`
		EventType     string                 `json:"event_type"`
		TargetID      int                    `json:"target_id"`
		TargetType    string                 `json:"target_type"`
		TargetDetails string                 `json:"target_details"`
	} `json:"properties"`
	Type string `json:"type"`
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received Request: %+v\n", req)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request processed successfully"))
}

func main() {
	http.HandleFunc("/post", postHandler)
	fmt.Println("Server is listening on port 1313...")
	if err := http.ListenAndServe(":1313", nil); err != nil {
		log.Fatal(err)
	}
}

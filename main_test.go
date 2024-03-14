package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostHandler(t *testing.T) {
	handler := http.HandlerFunc(postHandler)

	tt := []struct {
		name   string
		method string
		body   string
		status int
	}{
		{
			name:   "valid request",
			method: "POST",
			body: `{
				"id": "1",
				"author_id": -3,
				"entity_id": 29,
				"entity_type": "Project",
				"details": {
					"author_name": "deploy-key-name",
					"author_class": "DeployKey",
					"target_id": 29,
					"target_type": "Project",
					"target_details": "example-project",
					"custom_message": {
						"protocol": "ssh",
						"action": "git-upload-pack"
					},
					"ip_address": "127.0.0.1",
					"entity_path": "example-group/example-project"
				},
				"ip_address": "127.0.0.1",
				"author_name": "deploy-key-name",
				"entity_path": "example-group/example-project",
				"target_details": "example-project",
				"created_at": "2022-07-26T05:43:53.662Z",
				"target_type": "Project",
				"target_id": 29,
				"event_type": "repository_git_operation"
			}`,
			status: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "/post", bytes.NewBufferString(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.status)
			}
		})
	}
}

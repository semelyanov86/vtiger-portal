package e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	// Create a new HTTP request and response
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Call the ErrorResponse function with a test error message and status code
	ErrorResponse(w, req, http.StatusBadRequest, "Test error message")

	// Check that the response has the correct status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}

	// Check that the response has the correct content type
	contentType := w.Header().Get("Content-Type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type %q, but got %q", expectedContentType, contentType)
	}

	// Decode the response body into a map
	var responseBody map[string]interface{}
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&responseBody)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Check that the response body has the expected error message
	expectedErrorMessage := "Test error message"
	actualErrorMessage, ok := responseBody["error"].(string)
	if !ok || actualErrorMessage != expectedErrorMessage {
		t.Errorf("Expected error message %q, but got %q", expectedErrorMessage, actualErrorMessage)
	}
}

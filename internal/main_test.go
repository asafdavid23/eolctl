package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetAvailableProducts tests the GetAvailableProducts function.
func TestGetAvailableProducts(t *testing.T) {
	// Set up a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Simulate a JSON response
		_, err := w.Write([]byte(`{"products":[{"name":"Go","version":"1.20"},{"name":"Python","version":"3.11"}]}`))
		if err != nil {
			t.Fatalf("could not write response: %v", err)
		}
	}))
	defer server.Close()

	// Replace the URL in GetAvailableProducts with the mock server URL
	url = server.URL // You'll need to declare 'url' at package level to make this work.

	// Call the function
	body, err := GetAvailableProducts()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the response
	expectedBody := `{"products":[{"name":"Go","version":"1.20"},{"name":"Python","version":"3.11"}]}`
	if string(body) != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}
}

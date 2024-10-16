package helpers_test

import (
	"eolctl/internal" // Import your package
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockResponse = []byte(`[ "product1", "product2", "product3" ]`)

// Create a variable to hold the mock function
var GetAvailableProductsFunc func() ([]byte, error)

// Initialize the function to use either the mock or the real implementation
func init() {
	// Set the default function to the real one
	GetAvailableProductsFunc = helpers.GetAvailableProducts
}

// In your main code or tests, you can set this to the mock
func useMock() {
	GetAvailableProductsFunc = func() ([]byte, error) {
		return []byte(mockResponse), nil // Return the mock response
	}
}

func TestGetAvailableProducts(t *testing.T) {
	// Create a new test server with a mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Override the URL in the function (you may need to refactor your code to accept a baseURL for testing)
	useMock()

	data, err := GetAvailableProductsFunc()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := []byte(`[ "product1", "product2", "product3" ]`)
	if string(data) != string(expected) {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

//

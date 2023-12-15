package ray

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDo(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	// Set up the test options
	opts := Options{
		URL:         server.URL,
		Method:      http.MethodGet,
		Timeout:     5,
		Query:       map[string]string{"key": "value"},
		Body:        nil,
		Header:      map[string]string{"Content-Type": "application/json"},
		ContentType: "application/json",
	}

	// Call the function being tested
	response, err := Do(opts)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response body
	expectedBody := "Hello, World!"
	if string(response) != expectedBody {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedBody, string(response))
	}
}
func TestDoRetry(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	// Set up the test options
	opts := Options{
		URL:         server.URL,
		Method:      http.MethodGet,
		Timeout:     5,
		Query:       map[string]string{"key": "value"},
		Body:        nil,
		Header:      map[string]string{"Content-Type": "application/json"},
		ContentType: "application/json",
		RetryTimes:  3,
	}

	// Call the function being tested
	response, err := DoRetry(opts)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response body
	expectedBody := "Hello, World!"
	if string(response) != expectedBody {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedBody, string(response))
	}
}
func TestDoJSON(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()

	// Set up the test options
	opts := Options{
		URL:         server.URL,
		Method:      http.MethodGet,
		Timeout:     5,
		Query:       map[string]string{"key": "value"},
		Body:        nil,
		Header:      map[string]string{"Content-Type": "application/json"},
		ContentType: "application/json",
		RetryTimes:  3,
	}

	// Define the expected data structure
	expectedData := struct {
		Message string `json:"message"`
	}{}

	// Call the function being tested
	err := DoJSON(opts, &expectedData)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response data
	expectedMessage := "Hello, World!"
	if expectedData.Message != expectedMessage {
		t.Errorf("Unexpected response data. Expected: %s, Got: %s", expectedMessage, expectedData.Message)
	}
}

func TestProxy(t *testing.T) {

	SetDefaultProxy("http://127.0.0.1:10808")
	// Set up the test options
	opts := Options{
		URL:     "https://www.google.com",
		Method:  http.MethodGet,
		Timeout: 5,
		Body:    nil,
		// Proxy:   "http://127.0.0.1:10808",
	}
	// Call the function being tested
	response, err := Do(opts)
	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Check the response body
	fmt.Println(string(response))
}
func TestDoStream(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write([]byte("Hello, World!\n"))
		w.Write([]byte("Hello, World-1!\n"))
		w.Write([]byte("Hello, World-3!\n"))

		w.Write([]byte("done\n"))
		// time.Sleep(6 * time.Second)
	}))
	defer server.Close()

	// Set up the test options
	opts := Options{
		URL:     server.URL,
		Method:  http.MethodGet,
		Timeout: 5,
		Query:   nil,
		Body:    nil,
		Header:  nil,
	}

	err := DoStream(opts, func(r *bufio.Reader) error {
		for {
			line, err := r.ReadBytes('\n')
			if err == io.EOF {
				fmt.Println(string(line))
				return nil
			}
			if err != nil {
				return err
			}
			fmt.Println(string(line))
		}
	})

	if err != nil {
		t.Errorf("X,Unexpected error: %v", err)
	}

	fmt.Println("Done")

}

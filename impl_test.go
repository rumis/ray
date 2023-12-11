package ray

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	SetGlobalLogger(StdLogger)
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write([]byte("Hello, " + r.URL.Query().Get("hello")))
	}))
	defer server.Close()

	// Call the function being tested
	response, err := Get(context.Background(), server.URL+"?hello=World")

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response body
	expectedBody := "Hello, World"
	if string(response) != expectedBody {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedBody, string(response))
	}
}

func TestGetJson(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write([]byte(`{"key":"value"}`))
	}))
	defer server.Close()

	// Call the function being tested
	var response map[string]string
	err := GetJson(context.Background(), server.URL, &response)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response body
	expectedBody := map[string]string{"key": "value"}
	if response["key"] != expectedBody["key"] {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedBody, response)
	}
}

func TestPostForm(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.PostFormValue("key")
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write([]byte(v))
	}))
	defer server.Close()

	// Call the function being tested
	resp, err := PostForm(context.Background(), server.URL, map[string]string{"key": "x-value"})

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response body
	expectedBody := "x-value"
	if string(resp) != expectedBody {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedBody, resp)
	}
}

func TestPostRaw(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write(v)
	}))
	defer server.Close()

	// Call the function being tested
	resp, err := PostRaw(context.Background(), server.URL, map[string]string{"key": "x-value"})

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response body
	expectedBody := `{"key":"x-value"}`
	if string(resp) != expectedBody {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedBody, resp)
	}
}

func TestPostFormJson(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.PostFormValue("key")
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write([]byte(`{"key":"` + v + `"}`))
	}))
	defer server.Close()

	// Call the function being tested
	var resp map[string]string
	err := PostFormJson(context.Background(), server.URL, map[string]string{"key": "x-value"}, &resp)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response body
	expectedBody := "x-value"
	if resp["key"] != expectedBody {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedBody, resp)
	}
}

func TestPostRawJson(t *testing.T) {
	SetGlobalLogger(StdLogger)
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Set the response status code
		w.WriteHeader(http.StatusOK)
		// Set the response body
		w.Write(v)
	}))
	defer server.Close()

	// Call the function being tested
	var resp map[string]string
	err := PostRawJson(context.Background(), server.URL, map[string]string{"key": "x-value"}, &resp, map[string]string{"h1": "v1"}, map[string]string{"h2": "v2"})

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the response body
	expectedBody := "x-value"
	if resp["key"] != expectedBody {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedBody, resp)
	}
}

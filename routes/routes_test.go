package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("incorrect status code received: received %v, expected %v",
			status,
			http.StatusOK)
	}

	expected := "Hello Web!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: received %v, expected %v",
			rr.Body.String(),
			expected)
	}
}

func TestApiHandler(t *testing.T) {
    req, err:= http.NewRequest("GET", "/api", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(apiHandler)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
		t.Errorf("incorrect status code received: received %v, expected %v",
			status,
			http.StatusOK)
	}

	expected := "api"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: received %v, expected %v",
			rr.Body.String(),
			expected)
    }
}

package server

import (
	"encoding/json"
	"github.com/Ludusamo/cool-lang-features/database"
	"github.com/gorilla/mux"
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
	server := Server{nil, mux.NewRouter()}
	handler := http.HandlerFunc(server.homeHandler())
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
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server := Server{nil, mux.NewRouter()}
	handler := http.HandlerFunc(server.apiHandler())
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

func TestFeatureGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/feature", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server := Server{database.CreateDatabase(), mux.NewRouter()}
    feat, addErr := server.db.AddFeature("Test", "Desc")
    if addErr != nil {
        t.Fatal(addErr)
    }
	handler := http.HandlerFunc(server.featuresHandler())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("incorrect status code received: received %v, expected %v",
			status,
			http.StatusOK)
	}

	var features []database.Feature
	err = json.NewDecoder(rr.Body).Decode(&features)
	if err != nil {
		t.Fatal(err)
	}
	if len(features) != 1 {
		t.Errorf("expected %v features, received %v features", 1, len(features))
	}
	if features[0] != *feat {
		t.Errorf("expected %v, received %v", feat, features[0])
	}
}

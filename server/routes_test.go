package server

import (
	"bytes"
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

func TestFeaturesGet(t *testing.T) {
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

func TestFeatureGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/feature/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server := Server{database.CreateDatabase(), mux.NewRouter()}
	feat, addErr := server.db.AddFeature("Test", "Desc")
	if addErr != nil {
		t.Fatal(addErr)
	}
	server.router.HandleFunc("/api/feature/{id:[0-9]+}", server.featureHandler())
	server.router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("incorrect status code received: received %v, expected %v",
			status,
			http.StatusOK)
	}

	var retrievedFeat database.Feature
	err = json.NewDecoder(rr.Body).Decode(&retrievedFeat)
	if err != nil {
		t.Fatal(err)
	}
	if retrievedFeat != *feat {
		t.Errorf("expected %v, received %v", feat, retrievedFeat)
	}
}

func TestFeaturePost(t *testing.T) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(
		map[string]string{"Name": "test", "Description": "desc"})
	req, err := http.NewRequest("POST", "/api/feature", b)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server := Server{database.CreateDatabase(), mux.NewRouter()}
	handler := http.HandlerFunc(server.featuresHandler())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("incorrect status code received: received %v, expected %v",
			status,
			http.StatusCreated)
	}

	var newFeature database.Feature
	err = json.NewDecoder(rr.Body).Decode(&newFeature)
	if newFeature.Name != "test" {
		t.Errorf("expected name to be %v, received %v", "test", newFeature.Name)
	}
}

func TestFeatureDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/feature/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server := Server{database.CreateDatabase(), mux.NewRouter()}
	_, addErr := server.db.AddFeature("Test", "Desc")
	if addErr != nil {
		t.Fatal(addErr)
	}
	server.router.HandleFunc("/api/feature/{id:[0-9]+}", server.featureHandler())
	server.router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("incorrect status code received: received %v, expected %v",
			status,
			http.StatusOK)
	}

	if feat, err := server.db.GetFeature(0); err == nil {
		t.Errorf("feature still exists in the database: %v", feat)
	}
}

func TestFeaturePatch(t *testing.T) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(
		map[string]string{"Name": "New Name", "Description": "New Desc"})
	req, err := http.NewRequest("PATCH", "/api/feature/0", b)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server := Server{database.CreateDatabase(), mux.NewRouter()}
	feat, addErr := server.db.AddFeature("Test", "Desc")
	if addErr != nil {
		t.Fatal(addErr)
	}
	server.router.HandleFunc("/api/feature/{id:[0-9]+}", server.featureHandler())
	server.router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("incorrect status code received: received %v, expected %v",
			status,
			http.StatusOK)
	}

	var modifiedFeat database.Feature
	err = json.NewDecoder(rr.Body).Decode(&modifiedFeat)
	if err != nil {
		t.Fatal(err)
	}
	if modifiedFeat != *feat {
		t.Errorf("expected %v, received %v", feat, modifiedFeat)
	}
	if modifiedFeat.Name != "New Name" {
		t.Errorf("expected name to be %v, received %v",
			"New Name",
			modifiedFeat.Name)
	}
}

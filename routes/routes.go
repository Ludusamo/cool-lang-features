package routes

import (
	"io"
	"net/http"
)

func RegisterHandlers() {
    http.HandleFunc("/", homeHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "Hello Web!")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

}

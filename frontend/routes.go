package server

import (
	"cool-lang-features/rpc"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

/** Implement serve http function for the router so it conforms to the Handler
 * interface. Goes through all mappings and checks their regex expressions for
 * a match on the url.
 * @lhs router pointer
 * @param w the response writer to write out the resulting output
 * @param r the http request object
 */
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mapping := range router.routes {
		match := mapping.regex.MatchString(r.URL.String())
		if match {
			mapping.handler(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

/** Adds a route to the router's list of routes.
 * Turns the regexString into a regular expression to be matched with later.
 * @lhs router pointer
 * @param regexStr regular expression string to match with
 * @param handler associated HTTP handler
 */
func (router *Router) AddRoute(regexStr string, handler http.HandlerFunc) {
	regex := regexp.MustCompile(regexStr)
	router.routes = append(router.routes, RegexHandlerMapping{regex, handler})
}

/** Registers handlers for all the different routes on the server
 * @lhs server pointer
 */
func (s *Server) RegisterHandlers() {
	s.router.AddRoute("/api/feature/([0-9]+)", s.featureHandler())
	s.router.AddRoute("/api/feature", s.featuresHandler())
	s.router.AddRoute("/api", s.apiHandler())
	// Needs to be placed last. These are evaluated in the order they are added
	s.router.AddRoute("/?[A-Za-z\\s\\/+]",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/index.html")
		})
}

/** Handler for the /api route
 * @lhs server pointer
 * @return Handler function for router
 */
func (s *Server) apiHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, "api")
	}
}

/** Handler for /api/feature for adding and retrieving features
 * @lhs server pointer
 * @return Handler function for router
 */
func (s *Server) featuresHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			res := rpc.GetFeatures(s.backend)
			if res.Err != "" {
				log.Fatal(res.Err)
			}
			json.NewEncoder(w).Encode(res.Data)
		case "POST":
			var featurePost struct {
				Name        string
				Description string
			}
			json.NewDecoder(r.Body).Decode(&featurePost)
			res := rpc.PostFeature(
				s.backend,
				featurePost.Name,
				featurePost.Description)
			if res.Err != "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res.Err)
			} else {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(res.Data)
			}
		}
	}
}

/** Handler for /api/feature/{id} for modifying, deleting and retrieving
 * features based on id
 * @lhs server pointer
 * @return Handler function for router
 */
func (s *Server) featureHandler() http.HandlerFunc {
	regex := regexp.MustCompile("/api/feature/([0-9]+)")
	return func(w http.ResponseWriter, r *http.Request) {
		id, parseErr := strconv.Atoi(regex.FindStringSubmatch(r.URL.String())[1])
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(parseErr)
			return
		}
		switch r.Method {
		case "GET":
			res := rpc.GetFeature(s.backend, id)
			if res.Err != "" {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(res.Err)
			} else {
				json.NewEncoder(w).Encode(res.Data)
			}
		case "PATCH":
			var featurePatch struct {
				Name        string
				Description string
			}
			json.NewDecoder(r.Body).Decode(&featurePatch)
			res := rpc.PatchFeature(
				s.backend,
				id,
				featurePatch.Name,
				featurePatch.Description)
			if res.Err != "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res.Err)
			} else {
				json.NewEncoder(w).Encode(res.Data)
			}
		case "DELETE":
			rpc.DeleteFeature(s.backend, id)
			w.WriteHeader(http.StatusOK)
		}
	}
}

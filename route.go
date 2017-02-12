package jsonrest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/jsonapi"
)

type jSONRoute struct {
	base string
	obj  Storable
}

// JSONAPIHandler : takes a list of strings of the base handler
type JSONAPIHandler struct {
	routes   []jSONRoute
	mimeType string
}

// JSONRequest : Enables unit testing of routeToVerb
type JSONRequest struct {
	URL    *url.URL
	Method string
	Body   io.ReadCloser
}

// NewJSONRequest : create a JSONRequest from a  http.Request
func NewJSONRequest(r *http.Request) *JSONRequest {
	return &JSONRequest{URL: r.URL, Method: r.Method, Body: r.Body}
}

// NewJSONAPIHandler : Create an JSONAPIHandler was no routes installed
func NewJSONAPIHandler() *JSONAPIHandler {
	return &JSONAPIHandler{routes: make([]jSONRoute, 0),
		mimeType: `application/vnd.api+json`}
}

// AddRouteVerbHandler : pass a Create Show List Delete verb struct with the base name
func (jsapi *JSONAPIHandler) AddRouteVerbHandler(routeBase string, obj Storable) {
	jsapi.routes = append(jsapi.routes, jSONRoute{fmt.Sprintf("/%v", routeBase), obj})
}

func (jsapi *JSONAPIHandler) routeToVerb(w http.ResponseWriter, r *JSONRequest) {
	for _, route := range jsapi.routes {
		if strings.HasPrefix(r.URL.Path, route.base) {
			path := strings.Split(r.URL.Path, "/")
			var err error
			// Switch on the REST verb
			switch {
			case r.Method == http.MethodPut:
				route.obj = route.obj.New()
				jsonapiRuntime := jsonapi.NewRuntime()
				err = jsonapiRuntime.UnmarshalPayload(r.Body, route.obj)
				if err == nil {
					err = route.obj.Create()
					if err == nil {
						w.WriteHeader(http.StatusCreated)
						w.Header().Set("Content-Type", jsonapi.MediaType)
						err = jsonapiRuntime.MarshalOnePayload(w, route.obj)
					}
				}
			case r.Method == http.MethodDelete:
				if len(path) > 2 && len(path[2]) > 0 {
					// error handling !
					id, _ := strconv.Atoi(path[2])
					err = route.obj.Delete(IDtype(id))
				} else {
					http.Error(w, "ID required with DELETE command", http.StatusBadRequest)
					return
				}
			case r.Method == http.MethodGet && len(path) > 2 && len(path[2]) > 0:
				id, err := strconv.Atoi(path[2])
				jsonapiRuntime := jsonapi.NewRuntime()
				err = route.obj.Show(IDtype(id))
				if err == nil {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", jsonapi.MediaType)
					err = jsonapiRuntime.MarshalOnePayload(w, route.obj)
				}
			case r.Method == http.MethodGet:
				//err = verbs.List(w, r)
				log.Println("List not implemented")
				//		err = route.obj.List()
				http.Error(w, "List not implemented", http.StatusNotImplemented)
			default:
				http.Error(w, "Unsupported HTTP verb", http.StatusBadRequest)
			}
			// Handle the error
			switch e := err.(type) {
			case nil:
				break
			default:
				http.Error(w, e.Error(), http.StatusInternalServerError)
				log.Printf("500 Internal Error: %v", e)
			}
			return
		} // end if path startswith x
	} // end path tests
	http.Error(w, "Not found", http.StatusNotFound)
}

// overly trusting placeholder code
func authorized(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func (jsapi *JSONAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !authorized(w, r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	if r.Header.Get("Accept") != jsapi.mimeType {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}
	jsapi.routeToVerb(w, NewJSONRequest(r))
}

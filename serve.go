package jsonrest

import (
	"log"
	"net/http"
)

// Init configure application handler
func Init() {
	handler := NewJSONAPIHandler()
	handler.AddRouteVerbHandler("contact", &Contact{})
	handler.AddRouteVerbHandler("telephone", &Telephone{})
	handler.AddRouteVerbHandler("address", &AddressPhysical{})
	handler.AddRouteVerbHandler("email", &AddressEmail{})
	log.Fatal(http.ListenAndServe(":8080", handler))
}

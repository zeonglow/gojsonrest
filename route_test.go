package jsonrest

import (
	"net/http"
	"net/url"
	"testing"
)

type StubResponseWriter struct {
	HeaderCalled, WriteCalled, WriteHeaderCalled bool
	ResponseCode                                 int
}

func (stub *StubResponseWriter) Header() http.Header {
	stub.HeaderCalled = true
	return make(http.Header)
}

func (stub *StubResponseWriter) Write([]byte) (int, error) {
	stub.WriteCalled = true
	return 0, nil
}

func (stub *StubResponseWriter) WriteHeader(status int) {
	stub.WriteHeaderCalled = true
	stub.ResponseCode = status
}

func TestAddRoute(t *testing.T) {
	jsonAPI := NewJSONAPIHandler()
	if len(jsonAPI.routes) != 0 {
		t.Fatal("New router is not empty")
	}

	jsonAPI.AddRouteVerbHandler("woot", nil)
	if len(jsonAPI.routes) != 1 {
		t.Fatal("AddRoute failed")
	}
}

func TestBlankRouter(t *testing.T) {
	jsonAPI := NewJSONAPIHandler()
	//	jsonAPI.validHeaders = func(*http.Request) bool { return true }
	stub := &StubResponseWriter{}
	rootURL, _ := url.Parse("/")
	stubRequest := &JSONRequest{rootURL, http.MethodGet, nil}

	jsonAPI.routeToVerb(stub, stubRequest)
	if stub.ResponseCode != 404 {
		t.Fatalf("actual %d, expected 404", stub.ResponseCode)
	}
}

func stubHandler(http.ResponseWriter, *JSONRequest) {
}

type stubVerbHandler struct {
	createHit bool
	showHit   bool
	listHit   bool
	deleteHit bool
}

func newStubVerbHandler() *stubVerbHandler {
	return &stubVerbHandler{}
}

func (v *stubVerbHandler) Create() error {
	v.createHit = true
	return nil
}

func (v *stubVerbHandler) Show(id IDtype) error {
	v.showHit = true
	return nil
}
func (v *stubVerbHandler) List() error {
	v.listHit = true
	return nil
}
func (v *stubVerbHandler) Delete(id IDtype) error {
	v.deleteHit = true
	return nil
}

func (v *stubVerbHandler) New() Storable {
	return v
}

type stubReadCloser struct {
}

func (s *stubReadCloser) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (s *stubReadCloser) Close() error {
	return nil
}

func TestWootRouter(t *testing.T) {
	sv := newStubVerbHandler()
	jsonAPI := NewJSONAPIHandler()
	jsonAPI.AddRouteVerbHandler("woot", sv)
	stub := &StubResponseWriter{}
	rootURL, _ := url.Parse("/woot")
	stubRequest := &JSONRequest{rootURL, http.MethodPut, &stubReadCloser{}}
	jsonAPI.routeToVerb(stub, stubRequest)
	if !sv.createHit {
		t.Fatal("Create method not called")
	}
	idURL, _ := url.Parse("/woot/12345")

	stubRequest = &JSONRequest{idURL, http.MethodDelete, &stubReadCloser{}}
	jsonAPI.routeToVerb(stub, stubRequest)
	if !sv.deleteHit {
		t.Fatal("Delete method not called")
	}

	stubRequest = &JSONRequest{idURL, http.MethodGet, &stubReadCloser{}}
	jsonAPI.routeToVerb(stub, stubRequest)
	if !sv.showHit {
		t.Fatal("Show method not called")
	}
	/*
		stubRequest = &JSONRequest{rootURL, http.MethodGet, &stubReadCloser{}}
		jsonAPI.routeToVerb(stub, stubRequest)
		if !sv.listHit {
			t.Fatal("List method not called")
		}
	*/
	sv = newStubVerbHandler()
	stubRequest = &JSONRequest{rootURL, `Candle`, &stubReadCloser{}}
	jsonAPI.routeToVerb(stub, stubRequest)
	if sv.listHit {
		t.Fatal("List method called when it should not have been")
	}

	if stub.ResponseCode == 404 {
		t.Fatalf("Not found called too")
	}

}

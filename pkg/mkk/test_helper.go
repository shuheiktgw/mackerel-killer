package mkk

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

// setup sets up a test HTTP server along with Mkk
// that is configured to talk to that test server.
func setup() (m *Mkk, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", mux)

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	m = NewMkk("")
	u, _ := url.Parse(server.URL + "/")
	m.Client.BaseURL = u

	return m, mux, server.URL, server.Close
}

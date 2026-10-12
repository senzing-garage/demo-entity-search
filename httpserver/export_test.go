package httpserver

import (
	"context"
	"net/http"
)

// NewRootMux exposes newRootMux to the httpserver_test package.
func (httpServer *BasicHTTPServer) NewRootMux(ctx context.Context) (*http.ServeMux, []string) {
	return httpServer.newRootMux(ctx)
}

package httpserver

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"github.com/senzing-garage/demo-entity-search/entitysearchservice"
	"github.com/senzing-garage/go-helpers/wraperror"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicHTTPServer is the default implementation of the HttpServer interface.
type BasicHTTPServer struct {
	AllowedHostnames     []string
	Arguments            []string
	AvoidServing         bool
	Command              string
	ConnectionErrorLimit int
	EnableAll            bool
	EnableEntitySearch   bool
	HTMLTitle            string
	KeepalivePingTimeout int
	MaxBufferSizeBytes   int
	ReadHeaderTimeout    time.Duration
	ServerAddress        string
	ServerPort           int
	TtyOnly              bool
	URLRoutePrefix       string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Serve method serves the httpservice over HTTP.

Input
  - ctx: A context to control lifecycle.
*/

func (httpServer *BasicHTTPServer) Serve(ctx context.Context) error {
	var err error

	rootMux, userMessages := httpServer.newRootMux(ctx)

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessages = append(userMessages, fmt.Sprintf("Starting server on interface:port '%s'...\n", listenOnAddress))

	for userMessage := range userMessages {
		outputln(userMessage)
	}

	server := http.Server{
		Addr:              listenOnAddress,
		Handler:           addIncomingRequestLogging(rootMux),
		ReadHeaderTimeout: httpServer.ReadHeaderTimeout,
	}

	// Start a web browser.  Unless disabled.

	if !httpServer.TtyOnly {
		_ = browser.OpenURL(fmt.Sprintf("http://localhost:%d", httpServer.ServerPort))
	}

	if !httpServer.AvoidServing {
		err = server.ListenAndServe()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (httpServer *BasicHTTPServer) newRootMux(ctx context.Context) (*http.ServeMux, []string) {
	var userMessages []string

	rootMux := http.NewServeMux()

	// Add to root Mux.

	userMessages = append(userMessages, httpServer.addEntitySearchToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addStatcHTMLToMux(ctx, rootMux)...)

	return rootMux, userMessages
}

func (httpServer *BasicHTTPServer) addEntitySearchToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	var result []string

	// if httpServer.EnableAll || httpServer.EnableEntitySearch {
	entitySearchMux := httpServer.getEntitySearchMux(ctx)
	rootMux.Handle("/entity-search/", http.StripPrefix("/entity-search", entitySearchMux))

	result = append(result, fmt.Sprintf(
		"Serving EntitySearch at                   http://localhost:%d/%s\n",
		httpServer.ServerPort,
		"entity-search",
	))
	// }

	return result
}

func (httpServer *BasicHTTPServer) addStatcHTMLToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	result := []string{}

	_ = ctx

	rootDir, err := fs.Sub(static, "static/root")
	if err != nil {
		panic(err)
	}

	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	return result
}

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *BasicHTTPServer) getEntitySearchMux(ctx context.Context) *http.ServeMux {
	service := &entitysearchservice.BasicHTTPService{}

	return service.Handler(ctx)
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}

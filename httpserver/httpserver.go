package httpserver

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"github.com/senzing-garage/demo-entity-search/entitysearchservice"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicHTTPServer is the default implementation of the HttpServer interface.
type BasicHTTPServer struct {
	AllowedHostnames         []string
	Arguments                []string
	AvoidServing             bool
	Command                  string
	ConnectionErrorLimit     int
	EnableAll                bool
	EnableEntitySearch       bool
	EnableSenzingRestAPI     bool
	GrpcDialOptions          []grpc.DialOption
	GrpcTarget               string
	HTMLTitle                string
	KeepalivePingTimeout     int
	LogLevelName             string
	MaxBufferSizeBytes       int
	ObserverOrigin           string
	Observers                []observer.Observer
	OpenAPISpecificationSpec []byte
	ReadHeaderTimeout        time.Duration
	Settings                 string
	SenzingInstanceName      string
	SenzingVerboseLogging    int64
	ServerAddress            string
	ServerOptions            []senzingrestapi.ServerOption
	ServerPort               int
	TtyOnly                  bool
	URLRoutePrefix           string
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
	var (
		err          error
		userMessages []string
	)

	rootMux := http.NewServeMux()

	// Add to root Mux.

	userMessages = append(userMessages, httpServer.addAPIToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addReverseProxyToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addEntitySearchToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addStatcHTMLToMux(ctx, rootMux)...)

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

	return wraperror.Errorf(err, "httpserver.Serve error: %w", err)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (httpServer *BasicHTTPServer) addAPIToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	var result []string

	// if httpServer.EnableAll || httpServer.EnableSenzingRestAPI || httpServer.EnableEntitySearch {
	senzingAPIMux := httpServer.getSenzingAPIMux(ctx)
	rootMux.Handle("/api/", http.StripPrefix("/api", senzingAPIMux))

	result = append(result,
		fmt.Sprintf("Serving Senzing REST API at               http://localhost:%d/%s", httpServer.ServerPort, "api"))
	// }

	return result
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

func (httpServer *BasicHTTPServer) addReverseProxyToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	var result []string

	// if httpServer.EnableAll || httpServer.EnableSenzingRestAPI || httpServer.EnableEntitySearch {
	senzingAPIMux2 := httpServer.getSenzingAPI2Mux(ctx)
	rootMux.Handle("/entity-search/api/", http.StripPrefix("/entity-search/api", senzingAPIMux2))

	result = append(result, fmt.Sprintf(
		"Serving Senzing REST API Reverse Proxy at http://localhost:%d/%s",
		httpServer.ServerPort,
		"entity-search/api",
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

func (httpServer *BasicHTTPServer) getSenzingAPIMux(ctx context.Context) *senzingrestapi.Server {
	return httpServer.getSenzingAPIGenericMux(ctx, "/api")
}

func (httpServer *BasicHTTPServer) getSenzingAPI2Mux(ctx context.Context) *senzingrestapi.Server {
	return httpServer.getSenzingAPIGenericMux(ctx, "/entity-search/api")
}

func (httpServer *BasicHTTPServer) getSenzingAPIGenericMux(
	ctx context.Context,
	urlRoutePrefix string,
) *senzingrestapi.Server {
	_ = ctx
	service := &senzingrestservice.BasicSenzingRestService{
		GrpcDialOptions:          httpServer.GrpcDialOptions,
		GrpcTarget:               httpServer.GrpcTarget,
		LogLevelName:             httpServer.LogLevelName,
		ObserverOrigin:           httpServer.ObserverOrigin,
		Observers:                httpServer.Observers,
		Settings:                 httpServer.Settings,
		SenzingInstanceName:      httpServer.SenzingInstanceName,
		SenzingVerboseLogging:    httpServer.SenzingVerboseLogging,
		URLRoutePrefix:           urlRoutePrefix,
		OpenAPISpecificationSpec: httpServer.OpenAPISpecificationSpec,
	}

	srv, err := senzingrestapi.NewServer(service, httpServer.ServerOptions...)
	if err != nil {
		log.Fatal(err)
	}

	return srv
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}

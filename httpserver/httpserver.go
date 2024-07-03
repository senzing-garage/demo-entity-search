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

func (httpServer *BasicHTTPServer) getSenzingAPIGenericMux(ctx context.Context, urlRoutePrefix string) *senzingrestapi.Server {
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
	rootMux := http.NewServeMux()
	userMessage := ""

	// Enable Senzing HTTP REST API.

	// if httpServer.EnableAll || httpServer.EnableSenzingRestAPI || httpServer.EnableEntitySearch {
	senzingAPIMux := httpServer.getSenzingAPIMux(ctx)
	rootMux.Handle("/api/", http.StripPrefix("/api", senzingAPIMux))
	userMessage = fmt.Sprintf("%sServing Senzing REST API at               http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "api")
	// }

	// Enable Senzing HTTP REST API as reverse proxy.

	// if httpServer.EnableAll || httpServer.EnableSenzingRestAPI || httpServer.EnableEntitySearch {
	senzingAPIMux2 := httpServer.getSenzingAPI2Mux(ctx)
	rootMux.Handle("/entity-search/api/", http.StripPrefix("/entity-search/api", senzingAPIMux2))
	userMessage = fmt.Sprintf("%sServing Senzing REST API Reverse Proxy at http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "entity-search/api")
	// }

	// Enable EntitySearch.

	// if httpServer.EnableAll || httpServer.EnableEntitySearch {
	entitySearchMux := httpServer.getEntitySearchMux(ctx)
	rootMux.Handle("/entity-search/", http.StripPrefix("/entity-search", entitySearchMux))
	userMessage = fmt.Sprintf("%sServing EntitySearch at                   http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "entity-search")
	// }

	// Add route to static files.

	rootDir, err := fs.Sub(static, "static/root")
	if err != nil {
		panic(err)
	}
	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessage = fmt.Sprintf("%sStarting server on interface:port '%s'...\n", userMessage, listenOnAddress)
	fmt.Println(userMessage)
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
	return err
}

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
	"github.com/senzing/demo-entity-search/entitysearchservice"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-rest-api-service/senzingrestapi"
	"github.com/senzing/go-rest-api-service/senzingrestservice"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// HttpServerImpl is the default implementation of the HttpServer interface.
type HttpServerImpl struct {
	AllowedHostnames               []string
	Arguments                      []string
	Command                        string
	ConnectionErrorLimit           int
	EnableAll                      bool
	EnableEntitySearch             bool
	EnableSenzingRestAPI           bool
	GrpcDialOptions                []grpc.DialOption
	GrpcTarget                     string
	HtmlTitle                      string
	KeepalivePingTimeout           int
	LogLevelName                   string
	MaxBufferSizeBytes             int
	ObserverOrigin                 string
	Observers                      []observer.Observer
	OpenApiSpecificationSpec       []byte
	ReadHeaderTimeout              time.Duration
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
	ServerAddress                  string
	ServerOptions                  []senzingrestapi.ServerOption
	ServerPort                     int
	TtyOnly                        bool
	UrlRoutePrefix                 string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

func (httpServer *HttpServerImpl) getSenzingApiGenericMux(ctx context.Context, urlRoutePrefix string) *senzingrestapi.Server {
	service := &senzingrestservice.SenzingRestServiceImpl{
		GrpcDialOptions:                httpServer.GrpcDialOptions,
		GrpcTarget:                     httpServer.GrpcTarget,
		LogLevelName:                   httpServer.LogLevelName,
		ObserverOrigin:                 httpServer.ObserverOrigin,
		Observers:                      httpServer.Observers,
		SenzingEngineConfigurationJson: httpServer.SenzingEngineConfigurationJson,
		SenzingModuleName:              httpServer.SenzingModuleName,
		SenzingVerboseLogging:          httpServer.SenzingVerboseLogging,
		UrlRoutePrefix:                 urlRoutePrefix,
		OpenApiSpecificationSpec:       httpServer.OpenApiSpecificationSpec,
	}
	srv, err := senzingrestapi.NewServer(service, httpServer.ServerOptions...)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *HttpServerImpl) getEntitySearchMux(ctx context.Context) *http.ServeMux {
	service := &entitysearchservice.HttpServiceImpl{}
	return service.Handler(ctx)
}

func (httpServer *HttpServerImpl) getSenzingApiMux(ctx context.Context) *senzingrestapi.Server {
	return httpServer.getSenzingApiGenericMux(ctx, "/api")
}

func (httpServer *HttpServerImpl) getSenzingApi2Mux(ctx context.Context) *senzingrestapi.Server {
	return httpServer.getSenzingApiGenericMux(ctx, "/entity-search/api")
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Serve method serves the httpservice over HTTP.

Input
  - ctx: A context to control lifecycle.
*/

func (httpServer *HttpServerImpl) Serve(ctx context.Context) error {
	rootMux := http.NewServeMux()

	userMessage := ""

	// Enable Senzing HTTP REST API.

	if httpServer.EnableAll || httpServer.EnableSenzingRestAPI || httpServer.EnableEntitySearch {
		senzingApiMux := httpServer.getSenzingApiMux(ctx)
		rootMux.Handle("/api/", http.StripPrefix("/api", senzingApiMux))
		userMessage = fmt.Sprintf("%sServing Senzing REST API at               http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "api")
	}

	// Enable Senzing HTTP REST API as reverse proxy.

	if httpServer.EnableAll || httpServer.EnableSenzingRestAPI || httpServer.EnableEntitySearch {
		senzingApiMux := httpServer.getSenzingApi2Mux(ctx)
		rootMux.Handle("/entity-search/api/", http.StripPrefix("/entity-search/api", senzingApiMux))
		userMessage = fmt.Sprintf("%sServing Senzing REST API Reverse Proxy at http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "entity-search/api")
	}

	// Enable EntitySearch.

	if httpServer.EnableAll || httpServer.EnableEntitySearch {
		entitySearchMux := httpServer.getEntitySearchMux(ctx)
		rootMux.Handle("/entity-search/", http.StripPrefix("/entity-search", entitySearchMux))
		userMessage = fmt.Sprintf("%sServing EntitySearch at                   http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "entity-search")
	}

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

	return server.ListenAndServe()
}

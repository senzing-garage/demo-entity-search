package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/senzing-garage/demo-entity-search/httpserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicHTTPServer_Serve(test *testing.T) {
	test.Parallel()

	httpServer := &httpserver.BasicHTTPServer{
		AvoidServing: true,
		TtyOnly:      true,
	}
	err := httpServer.Serve(test.Context())
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

func TestBasicHTTPServer_newRootMux(test *testing.T) {
	test.Parallel()

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "root redirect page",
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "entity search UI",
			path:           "/entity-search/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "removed Senzing REST API",
			path:           "/api/heartbeat",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "removed Senzing REST API proxy",
			path:           "/entity-search/api/heartbeat",
			expectedStatus: http.StatusNotFound,
		},
	}

	httpServer := &httpserver.BasicHTTPServer{}
	rootMux, _ := httpServer.NewRootMux(test.Context())

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			request := httptest.NewRequestWithContext(test.Context(), http.MethodGet, testCase.path, nil)
			recorder := httptest.NewRecorder()
			rootMux.ServeHTTP(recorder, request)
			assert.Equal(test, testCase.expectedStatus, recorder.Code)
		})
	}
}

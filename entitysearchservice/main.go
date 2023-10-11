package entitysearchservice

import (
	"context"
	"net/http"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The HttpService interface...
type HttpService interface {
	Handler(ctx context.Context) *http.ServeMux
}

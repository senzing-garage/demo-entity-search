package entitysearchservice

import (
	"context"
	"net/http"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The HTTPService interface...
type HTTPService interface {
	Handler(ctx context.Context) *http.ServeMux
}

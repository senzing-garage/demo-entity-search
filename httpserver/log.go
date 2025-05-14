package httpserver

import (
	"net/http"

	"github.com/senzing-garage/demo-entity-search/internal/log"
)

// createRequestLog returns a logger with relevant request fields.
func createRequestLog(request *http.Request, additionalFields ...map[string]interface{}) log.Logger {
	fields := map[string]interface{}{}
	if len(additionalFields) > 0 {
		fields = additionalFields[0]
	}

	if request != nil {
		fields["host"] = request.Host
		fields["remote_addr"] = request.RemoteAddr
		fields["method"] = request.Method
		fields["protocol"] = request.Proto
		fields["path"] = request.URL.Path
		fields["request_url"] = request.URL.String()
		fields["user_agent"] = request.UserAgent()
		fields["cookies"] = request.Cookies()
	}

	return log.WithFields(fields)
}

package entitysearchservice

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicHTTPService is the default implementation of the HttpServer interface.
type BasicHTTPService struct {
	HTMLTitle      string
	URLRoutePrefix string
}

type TemplateVariables struct {
	HTMLTitle      string
	URLRoutePrefix string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

// createRequestLog returns a logger with relevant request fields
// func createRequestLog(r *http.Request, additionalFields ...map[string]interface{}) log.Logger {
// 	fields := map[string]interface{}{}
// 	if len(additionalFields) > 0 {
// 		fields = additionalFields[0]
// 	}
// 	if r != nil {
// 		fields["host"] = r.Host
// 		fields["remote_addr"] = r.RemoteAddr
// 		fields["method"] = r.Method
// 		fields["protocol"] = r.Proto
// 		fields["path"] = r.URL.Path
// 		fields["request_url"] = r.URL.String()
// 		fields["user_agent"] = r.UserAgent()
// 		fields["cookies"] = r.Cookies()
// 	}
// 	return log.WithFields(fields)
// }

// func getCreateLogger(connectionUUID string, r *http.Request) xtermjs.Logger {
// 	createRequestLog(r, map[string]interface{}{"connection_uuid": connectionUUID}).Infof(
// "created logger for connection '%s'", connectionUUID)
// 	return createRequestLog(nil, map[string]interface{}{"connection_uuid": connectionUUID})
// }

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// func (httpService *BasicHTTPService) populateStaticTemplate(responseWriter http.ResponseWriter,
// request *http.Request, filepath string, templateVariables TemplateVariables) {
// 	_ = request

// 	templateBytes, err := static.ReadFile(filepath)
// 	if err != nil {
// 		http.Error(responseWriter, http.StatusText(500), 500)
// 		return
// 	}

// 	templateParsed, err := template.New("HtmlTemplate").Parse(string(templateBytes))
// 	if err != nil {
// 		http.Error(responseWriter, http.StatusText(500), 500)
// 		return
// 	}

// 	err = templateParsed.Execute(responseWriter, templateVariables)
// 	if err != nil {
// 		http.Error(responseWriter, http.StatusText(500), 500)
// 		return
// 	}
// }

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Handler method...

Input
  - ctx: A context to control lifecycle.

Output
  - http.ServeMux...
*/

func (httpService *BasicHTTPService) Handler(ctx context.Context) *http.ServeMux {
	_ = ctx
	rootMux := http.NewServeMux()

	// Add route to static files.

	rootDir, err := fs.Sub(static, "static/root")
	if err != nil {
		panic(err)
	}

	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	return rootMux
}

package chizerolog

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"runtime/debug"
)

// Recover recovers any call to panic() made by a request handler or middleware.
// It also logs an error-levelled message using the global zerolog logger.
// This middleware should be registered first (or second if Logger is also used).
func Recover(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			scheme := "http"
			if request.TLS != nil {
				scheme = "https"
			}
			url := fmt.Sprintf("%s://%s%s", scheme, request.Host, request.RequestURI)

			if rec := recover(); rec != nil {
				log.Error().
					Str("proto", request.Proto).
					Str("method", request.Method).
					Str("route", url).
					Interface("recovered", rec).
					Bytes("stack", debug.Stack()).
					Msg("A request handler has panicked.")
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	}
	return http.HandlerFunc(fn)
}

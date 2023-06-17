package chizerolog

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const dataKey = "chzl_meta"

// Logger uses the global zerolog logger to log HTTP requests.
// Log messages are printed with the debug level.
// This middleware should be registered first.
func Logger(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), dataKey, make(map[string]any)))

		proxy := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

		start := time.Now()
		defer func() {
			end := time.Now()

			scheme := "http"
			if request.TLS != nil {
				scheme = "https"
			}
			url := fmt.Sprintf("%s://%s%s", scheme, request.Host, request.RequestURI)

			var err error
			data := request.Context().Value(dataKey)
			if data != nil {
				injErr, ok := data.(map[string]any)["err"]
				if ok {
					err = injErr.(error)
				}
			}

			if err == nil {
				log.Debug().
					Str("proto", request.Proto).
					Str("method", request.Method).
					Str("route", url).
					Str("client_address", request.RemoteAddr).
					Int("response_size", proxy.BytesWritten()).
					Str("elapsed", fmt.Sprintf("%s", end.Sub(start))).
					Int("status_code", proxy.Status()).
					Msg("An incoming request has been processed.")
			} else {
				log.Error().
					Err(err).
					Str("proto", request.Proto).
					Str("method", request.Method).
					Str("route", url).
					Str("client_address", request.RemoteAddr).
					Int("response_size", proxy.BytesWritten()).
					Str("elapsed", fmt.Sprintf("%s", end.Sub(start))).
					Int("status_code", proxy.Status()).
					Msg("An incoming request has been processed and resulted in an unexpected error.")
			}
		}()

		next.ServeHTTP(proxy, request)
	}
	return http.HandlerFunc(fn)
}

// InjectError injects the given error to a specific key so that Logger will log its occurrence later on in the request chain.
func InjectError(request *http.Request, err error) {
	data := request.Context().Value(dataKey)
	if data == nil {
		return
	}
	data.(map[string]any)["err"] = err
}

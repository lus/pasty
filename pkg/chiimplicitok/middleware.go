package chiimplicitok

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// Middleware sets the status code of a request to http.StatusOK if it was not set explicitly by any handler.
func Middleware(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		proxy := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

		defer func() {
			if proxy.Status() == 0 {
				proxy.WriteHeader(http.StatusOK)
			}
		}()

		next.ServeHTTP(writer, request)
	}
	return http.HandlerFunc(fn)
}

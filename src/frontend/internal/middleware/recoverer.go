package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func Recoverer(log *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					log.Error(rvr.(string))

					http.Error(w, "Unexpected error occurred. Please try again later", http.StatusBadRequest)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

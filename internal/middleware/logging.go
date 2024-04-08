package middleware

import (
	"fmt"
	"log"
	"net/http"
)

type wrappedWriter struct {
    http.ResponseWriter
    status int
}

func (w *wrappedWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}



func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        wr := &wrappedWriter{w, http.StatusOK}
        next.ServeHTTP(wr, r)

        ip := r.RemoteAddr
        msg := fmt.Sprintf("%s %d %s %s %s\n", ip, wr.status, r.Method, r.RequestURI, r.UserAgent())
        log.Print(msg)
    })
}

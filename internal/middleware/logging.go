package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	//"github.com/coltonmosier/api-v1/internal/database"
)

type wrappedWriter struct {
    http.ResponseWriter
    status int
}

func (w *wrappedWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}

//var LogDB = database.InitLoggingDatabase()


func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // NOTE: I have to do this here bc it was overwriting the Content-Type header in the response
        w.Header().Add("Content-Type", "application/json")
        wr := &wrappedWriter{w, http.StatusOK}
        next.ServeHTTP(wr, r)

        ip := strings.Split(r.RemoteAddr, ":")[0]
        msg := fmt.Sprintf("%s %d %s %s %s\n", ip, wr.status, r.Method, r.RequestURI, r.UserAgent())
        log.Print(msg)

        // NOTE: This is for logging to the db

        //ipnodots := strings.Replace(ip, ".", "", -1)
        //query := "INSERT INTO logs (ip, method, url) VALUES ($1, $2, $3)"
        //_, err := LogDB.Exec(query, ip, r.Method, r.RequestURI)
        //if err != nil {
        //    log.Fatal("Error inserting log into database")
        //}

    })
}

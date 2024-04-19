package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

// var LogDB = database.InitLoggingDatabase()
var LogFile = os.Getenv("LOG_FILE")

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: I have to do this here bc it was overwriting the Content-Type header in the response
		start := time.Now()
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		wr := &wrappedWriter{w, http.StatusOK}
		next.ServeHTTP(wr, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		ip := strings.Split(r.Header.Get("X-Real-IP"), ":")[0]
		if ip == "" {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		}
		msg := fmt.Sprintf("%s %d %s %s %v\n", ip, wr.status, r.Method, r.RequestURI, time.Since(start))
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

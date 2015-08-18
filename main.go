package main

import (
	"flag"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
	//"path/filepath"
	"time"
)

var PATH = flag.String("path", "public", "directory of web files")
var PORT = flag.String("port", "3000", "On which port the HTTP server listens")

func main() {

	flag.Parse()
	log.Printf("goServe Started, listening on PORT: %s", *PORT)
	requestHandlers := alice.New(loggingHandler, recoverHandler, imgHandler)
	http.Handle("/favicon.ico", requestHandlers.ThenFunc(favIcoHandler))
	http.Handle("/", requestHandlers.ThenFunc(indexHandler))
	http.ListenAndServe(":"+*PORT, nil)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, *PATH+"/index.html")
}

func favIcoHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, *PATH+"/favicon.ico")
}

func imgHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		reqUrl := req.URL.String()
		if _, err := os.Stat(*PATH + reqUrl); os.IsNotExist(err) {
			next.ServeHTTP(rw, req)
		} else {
			http.ServeFile(rw, req, *PATH+reqUrl)
		}
	}

	return http.HandlerFunc(fn)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		end := time.Now()
		log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), end.Sub(start))
	}

	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %+v", err)
				http.Error(rw, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

package main

import (
	"flag"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dir = flag.String("directory", "public/", "directory of web files")

func main() {

	// Os specifics. Signal handlers and flag parser
	flag.Parse()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()

	hostPort := "127.0.0.1:8080"
	log.Printf("goServe Started. Listening on %s", hostPort)

	requestHandlers := alice.New(loggingHandler, recoverHandler)
	http.Handle("/", requestHandlers.ThenFunc(indexHandler))
	http.Handle("/favicon.ico", requestHandlers.ThenFunc(favIcoHandler))
	http.ListenAndServe(hostPort, nil)

	<-done
	log.Println("goServe exiting")
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "public/index.html")
}

func favIcoHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "public/favicon.ico")
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

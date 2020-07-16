package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/michalswi/got-you.git/api"
	"github.com/michalswi/got-you.git/server"
)

const (
	apiPath = "/x"
)

var user string
var getPass string
var postPass string

func main() {

	logger := log.New(os.Stdout, "fake file server ", log.LstdFlags|log.Lshortfile)
	serverAddress := os.Getenv("SERVER_ADDR")

	lnxfile := "bin/msconfig"
	winfile := "bin/msconfig.exe"

	muxR := mux.NewRouter()
	r := muxR.PathPrefix(apiPath).Subrouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		checkIP(w, r, logger)
		http.ServeFile(w, r, "index.html")
	})

	r.HandleFunc("/dw", func(w http.ResponseWriter, r *http.Request) {
		checkIP(w, r, logger)
		ua := r.Header.Get("User-Agent")
		ua = strings.ToLower(ua)
		switch {
		case strings.Contains(ua, "linux"):
			logger.Printf("linux user agent: %s \n", ua)
			http.ServeFile(w, r, lnxfile)
		case strings.Contains(ua, "curl"):
			logger.Printf("linux/curl user agent: %s \n", ua)
			http.ServeFile(w, r, lnxfile)
		case strings.Contains(ua, "wget"):
			logger.Printf("linux/wget user agent: %s \n", ua)
			http.ServeFile(w, r, lnxfile)
		case strings.Contains(ua, "windows"):
			logger.Printf("windows user agent: %s \n", ua)
			http.ServeFile(w, r, winfile)
		}
	})

	r.HandleFunc("/hz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	r.Path("/get").Methods("GET").HandlerFunc(api.BasicAuth(logger, user, getPass, "like your t-shirt"))
	r.Path("/post").Methods("POST").HandlerFunc(api.BasicAuth(logger, user, postPass, "like your t-shirt"))

	srv := server.NewServer(r, serverAddress)

	go func() {
		logger.Printf("Starting server on port %s...\n", serverAddress)
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	gracefulShutdown(srv, logger)
}

func gracefulShutdown(srv *http.Server, logger *log.Logger) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	logger.Printf("Shutting down the server...\n")
	os.Exit(0)
}

func checkIP(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		logger.Printf("from '%s' - someone IP: %s \n", r.URL.Path, forwarded)
	}
	logger.Printf("from '%s' - someone IP: %s \n", r.URL.Path, r.RemoteAddr)
}

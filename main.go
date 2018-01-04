// The main package of Tangram product.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultAddress  = ":8080"
	shutdownTimeout = 5 * time.Second
	// application exit status
	successExitStatus            = 0
	errorStopintServerStatusCode = 1
)

var (
	version   = "development"
	build     = "undefined"
	buildDate = "unknown"
	isReady   = false
)

// Healthy check handler.
// Used to verify if the application is running.
// An application is healthy if its healthy status code is >= 200 && <400
func healthyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// Readiness check.
// Used to verify if application is ready to serve client request.
// An application is ready if status code is >= 200 && <400
func readyHandler(w http.ResponseWriter, r *http.Request) {
	if isReady {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "NO")
	}
}

func address() string {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = defaultAddress
	}
	return addr
}

// print application info
func printBanner() {
	log.Println("Tangram")
	log.Printf("  version:       %s\n", version)
	log.Printf("  build:         %s\n", build)
	log.Printf("  build date:    %s\n", buildDate)
	log.Printf("  starting time: %s\n", time.Now().Format(time.RFC3339))
}

func startHTTPServer() *http.Server {
	// configure HTTP server and register application status entrypoints
	server := &http.Server{Addr: address()}
	http.HandleFunc("/healthy", healthyHandler)
	http.HandleFunc("/ready", readyHandler)
	go func() {
		log.Printf("Listening on %s\n", address())
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	isReady = true
	return server
}

func gracefulShutdown(server *http.Server) {
	// deal with Ctrl+C (SIGTERM) and graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	isReady = false
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	log.Printf("Shutting down, with a timeout of %s.\n", shutdownTimeout)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error stopping http server. Error: %v\n", err)
		os.Exit(errorStopintServerStatusCode)
	} else {
		log.Println("Tangram server stoped")
		os.Exit(successExitStatus)
	}
}

func main() {
	printBanner()
	server := startHTTPServer()
	gracefulShutdown(server)
}

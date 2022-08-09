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

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/readiness", readinessHandler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/add", addHandler)
	srv := &http.Server{
		Handler: r,
		Addr:    ":8081",
	}

	// start the server
	go func() {
		log.Println("Starting the server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(srv)
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	log.Printf("received request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("Hello %s\n", name)))

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("server is healthy..."))
	w.WriteHeader(http.StatusOK)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	a := 5
	b := 6

	sum := a + b

	w.Write([]byte(fmt.Sprintf("%d", sum)))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("server is ready to serve..."))
	w.WriteHeader(http.StatusOK)
}

func waitForShutdown(srv *http.Server) {
	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-intChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down ....")
	os.Exit(0)
}

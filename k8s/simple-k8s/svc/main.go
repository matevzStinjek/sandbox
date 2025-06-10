package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ADDR = ":8080"
)

type Server struct {
	srv *http.Server
}

func NewServer() *Server {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ADDR,
		Handler: mux,
	}

	mux.HandleFunc("/work", cpuIntensiveHandler)
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/", rootHandler)

	return &Server{
		srv: srv,
	}
}

func (s *Server) Start() error {
	log.Printf("starting server on %s", ADDR)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{ "status": "ok" }`))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
	w.WriteHeader(200)
	w.Write([]byte(`ok`))
}

func cpuIntensiveHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	targetDuration := 100 * time.Millisecond

	iterations := 0
	data := []byte("initial data")

	// Hash until we've spent at least 100ms doing CPU work
	for time.Since(startTime) < targetDuration {
		hash := sha256.Sum256(data)
		data = hash[:]
		iterations++
	}

	actualDuration := time.Since(startTime)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Performed %d hashing operations in %v\n",
		iterations, actualDuration)))
}

func run(rootCtx context.Context) error {
	ctx, stop := signal.NotifyContext(rootCtx, syscall.SIGTERM, os.Interrupt)
	defer stop()

	srv := NewServer()

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.Start()
	}()

	select {
	case err := <-srvErr:
		return err
	case <-ctx.Done():
		shutdownCtx, shtCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shtCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(context.Background()); err != nil && err != http.ErrServerClosed {
		log.Println(err)
	}
}

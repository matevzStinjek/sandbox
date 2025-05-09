package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer() *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handleHealth)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	val := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	json, err := json.Marshal(val)
	if err != nil {
		return
	}

	w.WriteHeader(400)
	w.Write(json)
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

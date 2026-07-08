package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HuskySlava/go-email-relay/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func WriteResponse(w http.ResponseWriter, status int, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, 200, map[string]string{"status": "ok"})
	})

	return mux
}

func StartServer(cfg *config.HTTPConfig) error {
	mux := newRouter()
	securedMux := securityHeaders(mux)

	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      securedMux,
		ReadTimeout:  time.Duration(cfg.TimeoutInSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.TimeoutInSeconds) * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting server...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server failed to start: %w", err)
	case <-quit:
		log.Printf("Server graceful shutdown")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Printf("Server stopped")

	return nil
}

package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	port := "9000"
	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":" + port,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	slog.Debug("starting server", "addr", "http://localhost:"+port)

	// Wait for interruption.
	select {
	case err := <-srvErr:
		// Error when starting HTTP server.
		slog.Error("server start failed", "err", err)
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	if err := srv.Shutdown(context.Background()); err != nil {
		slog.Error("server shutdown failed", "err", err)
		return
	}
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
		return
	})
	return mux
	//handler := otelhttp.NewHandler(mux, "go-hello")
	//return handler
}

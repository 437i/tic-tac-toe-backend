package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func NewServer(router *mux.Router) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}

func StartServer(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return fmt.Errorf("listen on %s: %w", server.Addr, err)
			}

			log.Printf("server started on http://localhost%s\n", server.Addr)

			go func() {
				if err := server.Serve(listener); err != nil &&
					!errors.Is(err, http.ErrServerClosed) {
					log.Printf("server stopped: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("server stopping")
			return server.Shutdown(ctx)
		},
	})
}

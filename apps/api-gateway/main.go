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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sayanmondal31/api-gateway/cache"
	"github.com/sayanmondal31/api-gateway/config"
	gwMiddleware "github.com/sayanmondal31/api-gateway/middleware"
	"github.com/sayanmondal31/api-gateway/proxy"
)

func main() {
	cfg := config.Load()

	// initialize redis
	if err := cache.Init(cfg.RedisURL); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	fmt.Println("Connected to redis!")

	// Initialize gateway proxy
	gp, err := proxy.NewGatewayProxy(cfg.AuthSvcURL)

	if err != nil {
		log.Fatalf("Failed to initialize gateway proxy: %v", err)
	}

	// initialize router
	r := chi.NewRouter()

	r.Use(middleware.RequestID) //?
	r.Use(middleware.RealIP)    // ?
	r.Use(middleware.Logger)    // this logs start and end
	r.Use(middleware.Recoverer) // recovers from panics

	// API Gateway healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"OK","service":"api-gateway"}`))
	})

	// // wildcard route : forward everything else to proxy router
	// r.HandleFunc("/*", gp.RouteRequest)

	r.Group(func(r chi.Router) {
		r.Post("/auth/register", gp.RouteRequest)
		r.Post("/auth/login", gp.RouteRequest)
		r.Post("/auth//verify", gp.RouteRequest)
		r.Post("/auth/refresh", gp.RouteRequest)
	})

	r.Group(func(r chi.Router) {
		r.Use(gwMiddleware.Authenticate(cfg.JWTSecret))

		r.Get("/auth/me", gp.RouteRequest)
		r.Post("/auth/logout", gp.RouteRequest)
	})

	server := &http.Server{
		Addr:        ":" + cfg.Port,
		Handler:     r,
		ReadTimeout: 15 * time.Second,
	}

	serverCtx, serverCancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)

	// if there is any shutdown event from kube
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// if there is graceful shutdown
	go func() {
		<-sig
		serverStopCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-serverStopCtx.Done()
			if serverStopCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out.. forcing exit.")
			}

		}()

		err := server.Shutdown(serverStopCtx)
		if err != nil {
			log.Fatal(err)
		}

		serverCancel()

	}()

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server faild to start: %v", err)
	}

	<-serverCtx.Done()
	fmt.Println("API Gateway stopped gracefully.")
}

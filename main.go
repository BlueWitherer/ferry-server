package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"ferry-srv/log"

	_ "ferry-srv/api"

	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

var visitors = cache.New(15*time.Minute, 30*time.Minute)

func getClientIP(r *http.Request) string {
	if cf := r.Header.Get("CF-Connecting-IP"); cf != "" {
		return cf
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func getVisitor(ip string) *rate.Limiter {
	if val, found := visitors.Get(ip); found {
		return val.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(10, 30)
	visitors.Set(ip, limiter, cache.DefaultExpiration)
	return limiter
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Print("Starting server...")

	port := os.Getenv("WEB_PORT")
	if port == "" {
		log.Warn("WEB_PORT variable is not set")
		port = "3000"
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: rateLimitMiddleware(http.DefaultServeMux),
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered from panic: %v", r)
			}
		}()

		log.Print("change this text")
		log.Done("Server started successfully! Serving at 0.0.0.0:%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Startup failed: %s", err.Error())
			return
		}
	}()

	// shutdown sequence
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Warn("Shutting down server...")

	log.Trace("Deferring cancel...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Debug("Shutting down HTTP...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Shutdown error: %s", err.Error())
		srv.Close()
	}
	log.Info("HTTP stopped")

	log.Print("Server stopped")
	log.Shutdown()
}

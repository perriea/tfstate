// Copyright © 2018 Aurelien PERRIER <a.perrier89@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

// Run Web Server
func Run(enableUI bool) {
	var (
		server            *http.Server
		router, ui, apiv1 *mux.Router
		fs                http.Handler
	)

	logger = log.New(os.Stdout, "HTTP: ", log.LstdFlags)
	logger.Println("Server is starting...")

	// Router
	router = mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	// health
	router.Handle("/healthz", healthz())

	// Subrouter API
	apiv1 = router.PathPrefix("/api/v1").Subrouter()

	// Routes
	// State
	apiv1.HandleFunc("/state/{project}", head).Methods("HEAD")
	apiv1.HandleFunc("/state/{project}", getState).Methods("GET")
	apiv1.HandleFunc("/state/{project}", postState).Methods("POST")
	apiv1.HandleFunc("/state/{project}", deleteState).Methods("DELETE")
	// Lock
	apiv1.HandleFunc("/lock/{project}", head).Methods("HEAD")
	apiv1.HandleFunc("/lock/{project}", lockState).Methods("LOCK")
	apiv1.HandleFunc("/lock/{project}", unlockState).Methods("UNLOCK")

	// Static
	if enableUI {
		ui = router.PathPrefix("/ui").Subrouter()
		fs = http.FileServer(http.Dir("./static/"))
		ui.Handle("/static/", http.StripPrefix("static/", fs))
	}

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Logging
	server = &http.Server{
		Addr:         ":8000",
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Printf("- \"%s %s %d %s %s\" - %s", r.Proto, r.Method, http.StatusOK, r.URL.Path, r.UserAgent(), requestID)
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func BasicAuth(w http.ResponseWriter, r *http.Request, username, password, realm string) bool {

	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return false
	}

	return true
}

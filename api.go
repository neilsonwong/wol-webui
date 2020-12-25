package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// ListenAndServe starts up the chi router
func ListenAndServe() {

	router := chi.NewRouter()
	router.Use(corsMiddleware())

	router.Route("/api", func(r chi.Router) {
		r.Route("/device", func(r chi.Router) {
			r.Route("/{deviceID}", func(r chi.Router) {
				r.Use(deviceCtx)
				r.Post("/wake", handleWake)
				r.Post("/ping", handlePing)
			})
			r.Get("/search", handleSearch)
		})
	})

	err := http.ListenAndServe(":8084", router)
	log.Fatal(err)
}

func deviceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceID := chi.URLParam(r, "deviceID")
		device := &Device{
			id:   deviceID,
			name: "test-device",
			kind: COMPUTER,
			ip:   "",
			mac:  "",
		}

		// if err != nil {
		// 	http.Error(w, http.StatusText(404), 404)
		// 	return
		// }
		ctx := context.WithValue(r.Context(), "device", device)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleWake(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	device, ok := ctx.Value("device").(*Device)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// send the wake

	fmt.Fprintf(w, "waking %s", device.id)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	device, ok := ctx.Value("device").(*Device)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// send the ping

	fmt.Fprintf(w, "pinging %s", device.id)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	// send the search

	fmt.Fprintf(w, "searching")
}

func corsMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

package main

import (
	"context"
	"encoding/json"
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
			r.Post("/", handleAddDevice)
			r.Route("/{deviceID}", func(r chi.Router) {
				r.Use(deviceCtx)
				r.Get("/", handleGetDevice)
				r.Put("/", handleUpdateDevice)
				r.Post("/wake", handleWake)
				r.Post("/ping", handlePing)
			})
		})
		r.Get("/devices", handleList)
	})

	err := http.ListenAndServe(":8084", router)
	log.Fatal(err)
}

func deviceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceID := chi.URLParam(r, "deviceID")
		device := GetDevice(deviceID)
		if device == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
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
	WakeDevice(*device)

	fmt.Fprintf(w, "waking %s", device.ID)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	device, ok := ctx.Value("device").(*Device)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// send the ping

	fmt.Fprintf(w, "pinging %s", device.ID)
}

func handleGetDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	device, ok := ctx.Value("device").(*Device)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAddDevice(w http.ResponseWriter, r *http.Request) {
	var device Device

	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	AddNewDevice(device)

	fmt.Fprintf(w, "updating %s", device.ID)
}

func handleUpdateDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceID")
	var device Device

	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil || device.ID != deviceID {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	UpdateDevice(device)

	fmt.Fprintf(w, "updated %s", device.ID)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	// send the search
	w.Header().Set("Content-Type", "application/json")
	d := ListDevices()

	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

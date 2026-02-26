package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Ensure uploads directory exists
	os.MkdirAll("static/uploads", 0755)
	os.MkdirAll("static/admin", 0755)

	// Initialize the database
	InitDB("events.db")

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// --- Public API routes ---

	// GET/POST /api/events
	http.HandleFunc("/api/events", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/events" || r.URL.Path == "/api/events/" {
			// POST requires admin auth
			if r.Method == http.MethodPost {
				requireAuth(HandleEvents).ServeHTTP(w, r)
			} else {
				HandleEvents(w, r)
			}
			return
		}
		http.NotFound(w, r)
	}))

	// Routes with event ID: /api/events/{id}, /api/events/{id}/book, /api/events/{id}/bookings
	http.HandleFunc("/api/events/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasSuffix(path, "/book") || strings.HasSuffix(path, "/bookings") {
			HandleBookings(w, r)
			return
		}

		// DELETE requires auth
		if r.Method == http.MethodDelete {
			requireAuth(HandleEventByID).ServeHTTP(w, r)
			return
		}

		HandleEventByID(w, r)
	}))

	// --- Admin API routes ---
	http.HandleFunc("/api/admin/login", corsMiddleware(HandleAdminLogin))
	http.HandleFunc("/api/admin/logout", corsMiddleware(HandleAdminLogout))
	http.HandleFunc("/api/admin/check", corsMiddleware(HandleAdminCheck))
	http.HandleFunc("/api/admin/stats", corsMiddleware(requireAuth(HandleAdminStats)))
	http.HandleFunc("/api/admin/events", corsMiddleware(requireAuth(HandleAdminEvents)))

	// Serve static files (includes /uploads/ and /admin/)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

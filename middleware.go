package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

// Admin credentials (hardcoded for demo)
const (
	adminUsername = "gracy"
	adminPassword = "barbie"
)

// Session store (in-memory)
var (
	sessions   = make(map[string]bool)
	sessionsMu sync.RWMutex
)

// generateSessionToken creates a random hex token
func generateSessionToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// createSession stores a new session and returns the token
func createSession() string {
	token := generateSessionToken()
	sessionsMu.Lock()
	sessions[token] = true
	sessionsMu.Unlock()
	return token
}

// deleteSession removes a session
func deleteSession(token string) {
	sessionsMu.Lock()
	delete(sessions, token)
	sessionsMu.Unlock()
}

// isValidSession checks if a session token is valid
func isValidSession(token string) bool {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()
	return sessions[token]
}

// getSessionToken extracts the session token from cookies
func getSessionToken(r *http.Request) string {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// requireAuth is middleware that checks for a valid admin session
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getSessionToken(r)
		if token == "" || !isValidSession(token) {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next(w, r)
	}
}

// corsMiddleware adds CORS headers for local development
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// HandleAdminLogin handles POST /api/admin/login
func HandleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	if req.Username != adminUsername || req.Password != adminPassword {
		writeError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token := createSession()
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]string{
			"username": adminUsername,
		},
	})
}

// HandleAdminLogout handles POST /api/admin/logout
func HandleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := getSessionToken(r)
	if token != "" {
		deleteSession(token)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Logged out"})
}

// HandleAdminCheck handles GET /api/admin/check
func HandleAdminCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := getSessionToken(r)
	if token == "" || !isValidSession(token) {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]string{
			"username": adminUsername,
		},
	})
}

// HandleAdminStats handles GET /api/admin/stats
func HandleAdminStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	stats, err := GetDashboardStats()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch stats")
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: stats})
}

// HandleAdminEvents handles GET /api/admin/events (events with booking counts)
func HandleAdminEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	events, err := GetAllEventsWithBookingCount()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch events")
		return
	}
	if events == nil {
		events = []EventWithBookingCount{}
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: events})
}

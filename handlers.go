package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Allowed image extensions
var allowedImageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// generateID creates a random hex string for booking IDs
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, APIResponse{Success: false, Message: message})
}

// HandleEvents handles GET /api/events and POST /api/events
func HandleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetEvents(w, r)
	case http.MethodPost:
		handleCreateEvent(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleEventByID handles GET /api/events/:id and DELETE /api/events/:id
func HandleEventByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractEventID(r.URL.Path, "/api/events/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		event, err := GetEventByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Event not found")
			} else {
				writeError(w, http.StatusInternalServerError, "Failed to fetch event")
			}
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: event})

	case http.MethodDelete:
		// Delete requires auth — checked at route level
		if err := DeleteEvent(id); err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeError(w, http.StatusNotFound, "Event not found")
			} else {
				writeError(w, http.StatusInternalServerError, "Failed to delete event")
			}
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Event deleted"})

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleBookings handles POST /api/events/:id/book and GET /api/events/:id/bookings
func HandleBookings(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasSuffix(path, "/book") && r.Method == http.MethodPost {
		handleCreateBooking(w, r)
		return
	}

	if strings.HasSuffix(path, "/bookings") && r.Method == http.MethodGet {
		handleGetBookings(w, r)
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func handleGetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := GetAllEvents()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch events")
		return
	}
	if events == nil {
		events = []Event{}
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: events})
}

func handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (32 MB max)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid multipart form data")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))
	date := strings.TrimSpace(r.FormValue("date"))
	location := strings.TrimSpace(r.FormValue("location"))
	totalSeatsStr := strings.TrimSpace(r.FormValue("total_seats"))

	// Validation
	if name == "" {
		writeError(w, http.StatusBadRequest, "Event name is required")
		return
	}
	if description == "" {
		writeError(w, http.StatusBadRequest, "Event description is required")
		return
	}
	if date == "" {
		writeError(w, http.StatusBadRequest, "Event date is required")
		return
	}
	if location == "" {
		writeError(w, http.StatusBadRequest, "Event location is required")
		return
	}
	totalSeats, err := strconv.Atoi(totalSeatsStr)
	if err != nil || totalSeats <= 0 {
		writeError(w, http.StatusBadRequest, "Total seats must be a positive number")
		return
	}

	// Handle image uploads
	files := r.MultipartForm.File["images"]
	if len(files) < 5 {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Minimum 5 images required (got %d)", len(files)))
		return
	}
	if len(files) > 10 {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Maximum 10 images allowed (got %d)", len(files)))
		return
	}

	var savedFiles []string
	for _, fh := range files {
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		if !allowedImageExts[ext] {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid image format '%s'. Allowed: jpg, jpeg, png, webp", ext))
			return
		}

		// Generate unique filename
		uniqueName := generateID() + ext
		src, err := fh.Open()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to read uploaded file")
			return
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join("static", "uploads", uniqueName))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to save uploaded file")
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to save uploaded file")
			return
		}
		savedFiles = append(savedFiles, uniqueName)
	}

	req := CreateEventRequest{
		Name:        name,
		Description: description,
		Date:        date,
		Location:    location,
		TotalSeats:  totalSeats,
		Images:      strings.Join(savedFiles, ","),
	}

	event, err := CreateEvent(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Message: "Event created successfully", Data: event})
}

func handleCreateBooking(w http.ResponseWriter, r *http.Request) {
	// Extract event ID: /api/events/123/book
	path := strings.TrimSuffix(r.URL.Path, "/book")
	id, err := extractEventID(path, "/api/events/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	// Validation
	if strings.TrimSpace(req.AttendeeName) == "" {
		writeError(w, http.StatusBadRequest, "Attendee name is required")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		writeError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if req.NumPasses < 1 || req.NumPasses > 5 {
		writeError(w, http.StatusBadRequest, "Number of passes must be between 1 and 5")
		return
	}

	bookingID := generateID()
	booking, err := CreateBooking(id, req, bookingID)
	if err != nil {
		if strings.Contains(err.Error(), "not enough seats") {
			writeError(w, http.StatusConflict, err.Error())
		} else if strings.Contains(err.Error(), "event not found") {
			writeError(w, http.StatusNotFound, "Event not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to create booking")
		}
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Booking confirmed!",
		Data:    booking,
	})
}

func handleGetBookings(w http.ResponseWriter, r *http.Request) {
	// Extract event ID: /api/events/123/bookings
	path := strings.TrimSuffix(r.URL.Path, "/bookings")
	id, err := extractEventID(path, "/api/events/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	// Verify event exists
	if _, err := GetEventByID(id); err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "Event not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to fetch event")
		}
		return
	}

	bookings, err := GetBookingsByEventID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}
	if bookings == nil {
		bookings = []Booking{}
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: bookings})
}

// extractEventID extracts the integer event ID from a URL path
func extractEventID(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	// Remove any trailing slashes or extra path segments
	if idx := strings.Index(idStr, "/"); idx != -1 {
		idStr = idStr[:idx]
	}
	return strconv.Atoi(idStr)
}

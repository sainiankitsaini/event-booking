package main

// Event represents an event in the system
type Event struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Date           string   `json:"date"`
	Location       string   `json:"location"`
	TotalSeats     int      `json:"total_seats"`
	AvailableSeats int      `json:"available_seats"`
	Images         string   `json:"-"`
	ImageList      []string `json:"images"`
	CreatedAt      string   `json:"created_at"`
}

// Booking represents a booking for an event
type Booking struct {
	ID           string `json:"id"`
	EventID      int    `json:"event_id"`
	AttendeeName string `json:"attendee_name"`
	Email        string `json:"email"`
	NumPasses    int    `json:"num_passes"`
	CreatedAt    string `json:"created_at"`
}

// BookingRequest is the incoming JSON for creating a booking
type BookingRequest struct {
	AttendeeName string `json:"attendee_name"`
	Email        string `json:"email"`
	NumPasses    int    `json:"num_passes"`
}

// CreateEventRequest holds form fields parsed from multipart form
type CreateEventRequest struct {
	Name        string
	Description string
	Date        string
	Location    string
	TotalSeats  int
	Images      string // comma-separated filenames
}

// APIResponse is a generic JSON response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// AdminLoginRequest is the incoming JSON for admin login
type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DashboardStats holds summary stats for the admin dashboard
type DashboardStats struct {
	TotalEvents      int `json:"total_events"`
	TotalBookings    int `json:"total_bookings"`
	TotalSeatsSold   int `json:"total_seats_sold"`
}

// EventWithBookingCount extends Event with a booking count for admin views
type EventWithBookingCount struct {
	Event
	BookingCount int `json:"booking_count"`
}

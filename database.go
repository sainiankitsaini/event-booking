package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB opens the SQLite database and creates tables if they don't exist
func InitDB(dbPath string) {
	var err error
	db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTables()
	seedData()
	log.Println("Database initialized successfully")
}

func createTables() {
	eventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		date TEXT NOT NULL,
		location TEXT NOT NULL,
		total_seats INTEGER NOT NULL,
		available_seats INTEGER NOT NULL,
		images TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	bookingsTable := `
	CREATE TABLE IF NOT EXISTS bookings (
		id TEXT PRIMARY KEY,
		event_id INTEGER NOT NULL,
		attendee_name TEXT NOT NULL,
		email TEXT NOT NULL,
		num_passes INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (event_id) REFERENCES events(id)
	);`

	// Migrate: add images column if it doesn't exist (for existing DBs)
	migrateImages := `ALTER TABLE events ADD COLUMN images TEXT DEFAULT '';`

	if _, err := db.Exec(eventsTable); err != nil {
		log.Fatalf("Failed to create events table: %v", err)
	}
	if _, err := db.Exec(bookingsTable); err != nil {
		log.Fatalf("Failed to create bookings table: %v", err)
	}
	// Ignore error if column already exists
	db.Exec(migrateImages)
}

// downloadImage downloads an image from a URL and saves it to static/uploads/
func downloadImage(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create("static/uploads/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func seedData() {
	// Check if events already exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check events count: %v", err)
	}
	if count > 0 {
		return // Already seeded
	}

	log.Println("Seeding sample events...")

	type seedEvent struct {
		Name        string
		Description string
		Date        string
		Location    string
		TotalSeats  int
		ImageIDs    []int // picsum photo IDs
	}

	sampleEvents := []seedEvent{
		{
			Name:        "Go Developer Meetup 2026",
			Description: "Join fellow Go developers for a day of talks, workshops, and networking. Topics include concurrency patterns, new Go features, and building production-ready microservices. Snacks and drinks provided!",
			Date:        "2026-04-15",
			Location:    "Tech Hub Conference Center, San Francisco",
			TotalSeats:  100,
			ImageIDs:    []int{1, 2, 3, 4, 5},
		},
		{
			Name:        "AI & Machine Learning Summit",
			Description: "A comprehensive summit covering the latest advancements in AI and ML. Featuring keynote speakers from leading tech companies, hands-on labs with LLMs, and panel discussions on responsible AI development.",
			Date:        "2026-05-20",
			Location:    "Innovation Center, New York City",
			TotalSeats:  250,
			ImageIDs:    []int{10, 11, 12, 13, 14},
		},
		{
			Name:        "Open Source Contributor Day",
			Description: "Spend a day contributing to popular open source projects with guidance from maintainers. Perfect for beginners and experienced developers alike. Bring your laptop and enthusiasm!",
			Date:        "2026-06-10",
			Location:    "Community Library Hall, Austin, TX",
			TotalSeats:  50,
			ImageIDs:    []int{20, 21, 22, 23, 24},
		},
	}

	for i, e := range sampleEvents {
		var imageFiles []string
		for _, pid := range e.ImageIDs {
			fname := fmt.Sprintf("seed_%d_%d.jpg", i+1, pid)
			url := fmt.Sprintf("https://picsum.photos/id/%d/800/500", pid)
			if err := downloadImage(url, fname); err != nil {
				log.Printf("Warning: failed to download seed image %s: %v", fname, err)
				// Use placeholder path anyway so gallery is not empty
			}
			imageFiles = append(imageFiles, fname)
		}
		imagesStr := strings.Join(imageFiles, ",")

		_, err := db.Exec(
			"INSERT INTO events (name, description, date, location, total_seats, available_seats, images) VALUES (?, ?, ?, ?, ?, ?, ?)",
			e.Name, e.Description, e.Date, e.Location, e.TotalSeats, e.TotalSeats, imagesStr,
		)
		if err != nil {
			log.Printf("Failed to seed event '%s': %v", e.Name, err)
		}
	}
	fmt.Println("Seeded 3 sample events")
}

// parseEvent converts raw DB images string to ImageList slice
func parseEvent(e *Event) {
	if e.Images != "" {
		e.ImageList = strings.Split(e.Images, ",")
	} else {
		e.ImageList = []string{}
	}
}

// GetAllEvents returns all events from the database
func GetAllEvents() ([]Event, error) {
	rows, err := db.Query("SELECT id, name, description, date, location, total_seats, available_seats, images, created_at FROM events ORDER BY date ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Date, &e.Location, &e.TotalSeats, &e.AvailableSeats, &e.Images, &e.CreatedAt); err != nil {
			return nil, err
		}
		parseEvent(&e)
		events = append(events, e)
	}
	return events, rows.Err()
}

// GetEventByID returns a single event by its ID
func GetEventByID(id int) (*Event, error) {
	var e Event
	err := db.QueryRow(
		"SELECT id, name, description, date, location, total_seats, available_seats, images, created_at FROM events WHERE id = ?", id,
	).Scan(&e.ID, &e.Name, &e.Description, &e.Date, &e.Location, &e.TotalSeats, &e.AvailableSeats, &e.Images, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	parseEvent(&e)
	return &e, nil
}

// CreateEvent inserts a new event and returns it
func CreateEvent(req CreateEventRequest) (*Event, error) {
	result, err := db.Exec(
		"INSERT INTO events (name, description, date, location, total_seats, available_seats, images) VALUES (?, ?, ?, ?, ?, ?, ?)",
		req.Name, req.Description, req.Date, req.Location, req.TotalSeats, req.TotalSeats, req.Images,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return GetEventByID(int(id))
}

// DeleteEvent removes an event and its bookings
func DeleteEvent(id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM bookings WHERE event_id = ?", id); err != nil {
		return err
	}
	res, err := tx.Exec("DELETE FROM events WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("event not found")
	}
	return tx.Commit()
}

// CreateBooking inserts a new booking and decrements available seats
func CreateBooking(eventID int, req BookingRequest, bookingID string) (*Booking, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check available seats within the transaction
	var available int
	err = tx.QueryRow("SELECT available_seats FROM events WHERE id = ?", eventID).Scan(&available)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}
	if available < req.NumPasses {
		return nil, fmt.Errorf("not enough seats available (requested %d, available %d)", req.NumPasses, available)
	}

	// Decrement seats
	_, err = tx.Exec("UPDATE events SET available_seats = available_seats - ? WHERE id = ?", req.NumPasses, eventID)
	if err != nil {
		return nil, err
	}

	// Insert booking
	_, err = tx.Exec(
		"INSERT INTO bookings (id, event_id, attendee_name, email, num_passes) VALUES (?, ?, ?, ?, ?)",
		bookingID, eventID, req.AttendeeName, req.Email, req.NumPasses,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Booking{
		ID:           bookingID,
		EventID:      eventID,
		AttendeeName: req.AttendeeName,
		Email:        req.Email,
		NumPasses:    req.NumPasses,
	}, nil
}

// GetBookingsByEventID returns all bookings for a given event
func GetBookingsByEventID(eventID int) ([]Booking, error) {
	rows, err := db.Query(
		"SELECT id, event_id, attendee_name, email, num_passes, created_at FROM bookings WHERE event_id = ? ORDER BY created_at DESC", eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.ID, &b.EventID, &b.AttendeeName, &b.Email, &b.NumPasses, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

// GetAllEventsWithBookingCount returns events with their booking counts for admin
func GetAllEventsWithBookingCount() ([]EventWithBookingCount, error) {
	rows, err := db.Query(`
		SELECT e.id, e.name, e.description, e.date, e.location, e.total_seats, e.available_seats, e.images, e.created_at,
		       COALESCE((SELECT COUNT(*) FROM bookings b WHERE b.event_id = e.id), 0) as booking_count
		FROM events e ORDER BY e.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []EventWithBookingCount
	for rows.Next() {
		var ewb EventWithBookingCount
		if err := rows.Scan(&ewb.ID, &ewb.Name, &ewb.Description, &ewb.Date, &ewb.Location,
			&ewb.TotalSeats, &ewb.AvailableSeats, &ewb.Images, &ewb.CreatedAt, &ewb.BookingCount); err != nil {
			return nil, err
		}
		parseEvent(&ewb.Event)
		results = append(results, ewb)
	}
	return results, rows.Err()
}

// GetDashboardStats returns summary statistics
func GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats
	err := db.QueryRow("SELECT COUNT(*) FROM events").Scan(&stats.TotalEvents)
	if err != nil {
		return nil, err
	}
	err = db.QueryRow("SELECT COUNT(*) FROM bookings").Scan(&stats.TotalBookings)
	if err != nil {
		return nil, err
	}
	err = db.QueryRow("SELECT COALESCE(SUM(num_passes), 0) FROM bookings").Scan(&stats.TotalSeatsSold)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

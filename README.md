# EventPass — Event Booking

A feature-rich event booking web application built with **Go** (backend) and **vanilla HTML/CSS/JS** (frontend), using **SQLite** as the database. Features a dark glassmorphism UI, admin panel with authentication, image galleries, and more.

## Features

- **Browse Events** — Public event listing with hero images, animated cards, and seat badges
- **Book Passes** — Reserve 1–5 passes with confetti confirmation and booking IDs
- **Image Galleries** — Each event supports 5–10 uploaded images with a featured image viewer
- **Admin Panel** — Session-based login, dashboard with stats, event management, CSV export
- **Toast Notifications** — Reusable toast system for success/error/info messages
- **Fancy UI** — Dark theme, glassmorphism cards, animated counters, floating particles
- **Seed Data** — 3 sample events with placeholder images auto-created on first run

## Admin Login

- **Username:** `gracy`
- **Password:** `barbie`
- **URL:** http://localhost:8080/admin/login.html

## Project Structure

```
event-booking/
├── main.go              # Server entry point and routing
├── database.go          # SQLite setup, queries, seed data
├── handlers.go          # Public API handlers (events, bookings)
├── handlers_admin.go    # Admin API handlers (login, stats, events)
├── middleware.go         # Auth middleware, session management, CORS
├── models.go            # Data structures
├── go.mod / go.sum
├── README.md
└── static/
    ├── index.html       # Public: event listing with hero
    ├── event.html       # Public: event detail + booking form
    ├── style.css        # Global styles (dark theme)
    ├── admin/
    │   ├── login.html       # Admin login page
    │   ├── dashboard.html   # Admin dashboard with stats + event table
    │   ├── create.html      # Create event with image upload
    │   └── event.html       # Admin event detail + bookings table
    ├── js/
    │   ├── toast.js     # Reusable toast notification system
    │   ├── gallery.js   # Image gallery component
    │   └── confetti.js  # Confetti animation wrapper (CDN)
    └── uploads/         # Auto-created, stores uploaded images
```

## Prerequisites

- **Go 1.18+** installed
- **GCC** available (required by `go-sqlite3` CGo dependency)
- Internet connection on first run (downloads seed images from picsum.photos)

## How to Run

```bash
# Install dependencies
go mod tidy

# Run the application
go run .
```

The app will be available at **http://localhost:8080**

Set a custom port with `PORT=3000 go run .`

## API Endpoints

### Public (no auth required)

| Method | Endpoint                    | Description              |
|--------|-----------------------------|--------------------------|
| GET    | `/api/events`               | List all events          |
| GET    | `/api/events/:id`           | Get single event details |
| POST   | `/api/events/:id/book`      | Book passes for an event |
| GET    | `/api/events/:id/bookings`  | List bookings for event  |

### Admin (auth required)

| Method | Endpoint                    | Description                  |
|--------|-----------------------------|------------------------------|
| POST   | `/api/admin/login`          | Login (sets session cookie)  |
| POST   | `/api/admin/logout`         | Logout (clears cookie)       |
| GET    | `/api/admin/check`          | Check auth status            |
| GET    | `/api/admin/stats`          | Dashboard statistics         |
| GET    | `/api/admin/events`         | Events with booking counts   |
| POST   | `/api/events`               | Create event (multipart)     |
| DELETE | `/api/events/:id`           | Delete event + bookings      |

## Deployment

### Build a Single Binary

```bash
go build -o event-booking .
```

### Deploy to a Server

```bash
# Copy binary + static folder to your server
scp -r event-booking static/ user@yourserver:/app/

# Run on server
cd /app && ./event-booking
```

### Run as Background Service

```bash
# Using screen
screen -S eventbooking
PORT=8080 ./event-booking
# Ctrl+A, D to detach
```

### Cloud Deployment

Works on **Railway**, **Fly.io**, or any VPS. Ensure GCC is available for the SQLite CGo driver.

## Tech Stack

- **Backend:** Go standard library (`net/http`)
- **Database:** SQLite via `github.com/mattn/go-sqlite3`
- **Frontend:** Vanilla HTML, CSS, JavaScript (no frameworks)
- **Fonts:** Inter + Playfair Display (Google Fonts)
- **Confetti:** canvas-confetti via CDN

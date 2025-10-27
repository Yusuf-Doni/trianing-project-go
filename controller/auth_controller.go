package controller

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
)

// LoginController handles login page display and authentication
func LoginController(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Check if user is already logged in
			if IsLoggedIn(r) {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}

			fp := filepath.Join("view", "login.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				log.Printf("Error parsing login template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				log.Printf("Error executing login template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else if r.Method == "POST" {
			// Handle login form submission
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username == "" || password == "" {
				http.Error(w, "Username and password are required", http.StatusBadRequest)
				return
			}

			// Authenticate user
			user, err := authenticateUser(db, username, password)
			if err != nil {
				log.Printf("Authentication error: %v", err)
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			// Create session
			sessionID, err := createSession(db, user.ID, user.Username)
			if err != nil {
				log.Printf("Session creation error: %v", err)
				http.Error(w, "Failed to create session", http.StatusInternalServerError)
				return
			}

			// Set session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				Secure:   false, // Set to true in production with HTTPS
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(24 * time.Hour),
			})

			log.Printf("User %s logged in successfully", username)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// LogoutController handles user logout
func LogoutController(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Get session ID from cookie
			cookie, err := r.Cookie("session_id")
			if err == nil {
				// Delete session from database
				deleteSession(db, cookie.Value)
			}

			// Clear session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Unix(0, 0),
			})

			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// RegisterController handles user registration
func RegisterController(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

			if IsLoggedIn(r) {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}

			fp := filepath.Join("view", "register.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				log.Printf("Error parsing register template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				log.Printf("Error executing register template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else if r.Method == "POST" {
			// Handle registration form submission
			username := r.FormValue("username")
			password := r.FormValue("password")
			email := r.FormValue("email")

			if username == "" || password == "" || email == "" {
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return
			}

			// Check if user already exists
			exists, err := userExists(db, username, email)
			if err != nil {
				log.Printf("Error checking user existence: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if exists {
				http.Error(w, "Username or email already exists", http.StatusConflict)
				return
			}

			// Create new user
			userID, err := createUser(db, username, password, email)
			if err != nil {
				log.Printf("Error creating user: %v", err)
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}

			log.Printf("New user registered: %s (ID: %d)", username, userID)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// authenticateUser validates user credentials
func authenticateUser(db *sql.DB, username, password string) (*model.User, error) {
	query := `
		SELECT id, username, password, email, role 
		FROM users 
		WHERE username = $1 AND password = $2
	`

	var user model.User
	err := db.QueryRow(query, username, password).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Role,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

// createSession creates a new user session
func createSession(db *sql.DB, userID int, username string) (string, error) {
	// Generate random session ID
	sessionID := generateSessionID()

	// Create session in database
	query := `
		INSERT INTO sessions (id, user_id, username, created_at, expires_at) 
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	_, err := db.Exec(query, sessionID, userID, username, now, expiresAt)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// deleteSession removes a session from the database
func deleteSession(db *sql.DB, sessionID string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := db.Exec(query, sessionID)
	return err
}

// userExists checks if a user with given username or email already exists
func userExists(db *sql.DB, username, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2`
	var count int
	err := db.QueryRow(query, username, email).Scan(&count)
	return count > 0, err
}

// createUser creates a new user in the database
func createUser(db *sql.DB, username, password, email string) (int, error) {
	query := `
		INSERT INTO users (username, password, email, role) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`

	var userID int
	err := db.QueryRow(query, username, password, email, "user").Scan(&userID)
	return userID, err
}

// generateSessionID generates a random session ID
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// IsLoggedIn checks if the user is currently logged in
func IsLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}

	// You would typically validate the session in the database here
	// For simplicity, we'll just check if the cookie exists
	return cookie.Value != ""
}

// GetCurrentUser returns the current logged-in user
func GetCurrentUser(db *sql.DB, r *http.Request) (*model.User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, fmt.Errorf("no session found")
	}

	query := `
		SELECT u.id, u.username, u.email, u.role 
		FROM users u 
		JOIN sessions s ON u.id = s.user_id 
		WHERE s.id = $1 AND s.expires_at > $2
	`

	var user model.User
	err = db.QueryRow(query, cookie.Value, time.Now()).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	return &user, nil
}

// RequireAuth middleware to protect routes
func RequireAuth(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate session in database
		_, err := GetCurrentUser(db, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}

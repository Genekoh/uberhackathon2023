package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const (
	port = ":8080"
	dsn  = "./database.db"
)

var (
	sessionManager *scs.SessionManager
	db             *sql.DB
)

type User struct {
	Id           int64
	UserName     string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash []byte
	Salary       int64
	// CreatedAt time.Time
	// UpdatedAt time.Time
}

type LoginCredentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func main() {
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer conn.Close()
	conn.SetConnMaxLifetime(0)
	conn.SetMaxIdleConns(50)
	conn.SetMaxOpenConns(50)
	db = conn
	log.Println("Connected to database")

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}
	log.Println("Pinged database")

	sessionManager = scs.New()
	store := sqlite3store.New(db)
	// sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Lifetime = 10 * time.Minute
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Store = store
	defer store.StopCleanup()

	// rows, err := db.Query("SELECT firstName FROM users")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var firstName string
	// for rows.Next() {
	// 	err = rows.Scan(&firstName)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("row is", firstName)
	// }
	// defer rows.Close()

	router := chi.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return sessionManager.LoadAndSave(h)
	})

	router.Get("/get", func(w http.ResponseWriter, r *http.Request) {
		u := sessionManager.GetString(r.Context(), "username")
		fmt.Println("heres the cookie: " + u)
		fmt.Println(r.Cookies())
		if u == "" {
			io.WriteString(w, "not authorized"+"\n"+sessionManager.GetString(r.Context(), "mykey"))
		} else {
			io.WriteString(w, "authorized as:\t"+u+"\n"+sessionManager.GetString(r.Context(), "mykey"))
		}
	})

	router.Get("/signin", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)

		// queries for user in database
		user, err := getUserByEmail(r.Context(), "johndoe@gmail.com")
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			enc.Encode("{'ok': false}")
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode("{'ok': false}")
			return
		}

		sessionManager.Put(r.Context(), "username", user.UserName)
		enc.Encode("{'ok': true}")
	})

	router.Get("/set/{smth}", func(w http.ResponseWriter, r *http.Request) {
		smth := chi.URLParam(r, "smth")
		sessionManager.Put(r.Context(), "mykey", smth)

		res := fmt.Sprintf("Put %v into session\n", smth)
		io.WriteString(w, res)
	})

	router.Route("/accounts", func(r chi.Router) {
		r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
			enc := json.NewEncoder(w)

			// decode body and check if json has required data
			var credentials LoginCredentials
			if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
				fmt.Println("1:", err)
				w.WriteHeader(http.StatusBadRequest)
				enc.Encode("{'ok': false}")
				return
			} else if credentials.Email == "" || credentials.Password == "" || len([]byte(credentials.Password)) > 72 {
				fmt.Println("2:", err)
				w.WriteHeader(http.StatusBadRequest)
				enc.Encode("{'ok': false}")
				return
			}

			// queries for user in database
			user, err := getUserByEmail(r.Context(), credentials.Email)
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				enc.Encode("{'ok': false}")
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				enc.Encode("{'ok': false}")
				return
			}

			// compare password and hash
			err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(credentials.Password))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				enc.Encode("{'ok': false}")
				return
			}

			sessionManager.Put(r.Context(), "username", user.UserName)
			enc.Encode("{'ok': true}")
		})

		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {})

	})

	server := http.Server{
		Addr:    port,
		Handler: router,
	}
	log.Printf("Starting server on port %s\n", port)
	server.ListenAndServe()
}

// get UserByUsername finds a User row from database and returns the object
func getUserByUsername(ctx context.Context, userName string) (*User, error) {
	var (
		username, firstName, lastName, email string
		passwordHash                         []byte
		userId, salary                       int64
	)
	row := db.QueryRowContext(ctx, "select rowid, * from users where userName = ?", userName)
	if err := row.Scan(&userId, &username, &firstName, &lastName, &email, &passwordHash, &salary); err != nil {
		return nil, err
	}

	return &User{userId, username, firstName, lastName, email, passwordHash, salary}, nil
}

func getUserByEmail(ctx context.Context, emailQuery string) (*User, error) {
	var (
		username, firstName, lastName, email string
		passwordHash                         []byte
		userId, salary                       int64
	)
	row := db.QueryRowContext(ctx, "select rowid, * from users where email = ?", emailQuery)
	if err := row.Scan(&userId, &username, &firstName, &lastName, &email, &passwordHash, &salary); err != nil {
		return nil, err
	}

	return &User{userId, username, firstName, lastName, email, passwordHash, salary}, nil
}

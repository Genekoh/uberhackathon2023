package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
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
	Id        int64
	UserName  string
	FirstName string
	LastName  string
	Email     string
	Password  []byte
	Salary    int64
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
	sessionManager.Lifetime = 5 * time.Second
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
		fmt.Println(sessionManager.GetString(r.Context(), "mykey"))

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, "select rowid, * from users")
		defer rows.Close()

		var (
			username, firstName, lastName, email string
			password                             []byte
			userId, salary                       int64
		)
		res := []User{}
		for rows.Next() {
			err = rows.Scan(&userId, &username, &firstName, &lastName, &email, &password, &salary)
			if err != nil {
				log.Fatal(err)
			}

			newUser := User{userId, username, firstName, lastName, email, password, salary}
			res = append(res, newUser)
		}

		enc := json.NewEncoder(w)
		enc.Encode(res)
	})

	router.Get("/set/{smth}", func(w http.ResponseWriter, r *http.Request) {
		smth := chi.URLParam(r, "smth")
		sessionManager.Put(r.Context(), "mykey", smth)

		res := fmt.Sprintf("Put %v into session\n", smth)
		w.Write([]byte(res))
	})

	router.Route("/accounts", func(r chi.Router) {
		r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
			enc := json.NewEncoder(w)

			var credentials LoginCredentials
			if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
				fmt.Println("1:", err)
				w.WriteHeader(http.StatusBadRequest)
				enc.Encode("{'ok': false,}")
				return
			}
			if credentials.Email == "" || credentials.Password == "" {
				fmt.Println("2:", err)
				w.WriteHeader(http.StatusBadRequest)
				enc.Encode("{'ok': false,}")
				return
			}

			user, err := getUserByEmail(r.Context(), credentials.Email)
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				enc.Encode("{'ok': false,}")
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				enc.Encode("{'ok': false,}")
				return
			}

			fmt.Printf("%+v\n", user)

			enc.Encode("{'ok': true,}")

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
		password                             []byte
		userId, salary                       int64
	)
	row := db.QueryRowContext(ctx, "select rowid, * from users where userName = ?", userName)
	if err := row.Scan(&userId, &username, &firstName, &lastName, &email, &password, &salary); err != nil {
		return nil, err
	}

	return &User{userId, username, firstName, lastName, email, password, salary}, nil
}

func getUserByEmail(ctx context.Context, emailQuery string) (*User, error) {
	var (
		username, firstName, lastName, email string
		password                             []byte
		userId, salary                       int64
	)
	row := db.QueryRowContext(ctx, "select rowid, * from users where email = ?", emailQuery)
	if err := row.Scan(&userId, &username, &firstName, &lastName, &email, &password, &salary); err != nil {
		return nil, err
	}

	return &User{userId, username, firstName, lastName, email, password, salary}, nil
}

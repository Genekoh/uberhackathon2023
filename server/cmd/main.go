package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

const (
	port = ":8080"
	dsn  = "./database.db"
)

var (
	sessionManager *scs.SessionManager
	db             *sql.DB
)

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
	// c.execute("PRAGMA foreign_keys = ON")

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatalf("failed to turn foreign_keys on: %v", err)
	}

	sessionManager = scs.New()
	store := sqlite3store.New(db)
	// sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Lifetime = 10 * time.Minute
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Store = store
	defer store.StopCleanup()

	router := chi.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return sessionManager.LoadAndSave(h)
	})

	router.Route("/accounts", func(r chi.Router) {
		r.Post("/signin", PostSignin)
		r.Post("/signup", PostSignup)
		r.Put("/update-salary", PostUpdateSalary)
	})

	router.Post("/book-ride", PostBookRide)

	// router.Handle("/listen-carpool", websocket.Handler(WsListenCarpool))

	carpoolsCleanup(context.Background(), carpoolTimeout)

	server := http.Server{
		Addr:    port,
		Handler: router,
	}
	log.Printf("Starting server on port %s\n", port)
	server.ListenAndServe()
}

func carpoolsCleanup(ctx context.Context, d time.Duration) {
	ticker := time.NewTicker(d)

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				log.Println("cleaning up carpools")
				_, err := db.Exec("DELETE FROM carpools WHERE ? > expiresAt", time.Now().Unix())
				if err != nil {
					log.Printf("carpools cleanup error: %v\n", err)
				}
			}
		}
	}()
}

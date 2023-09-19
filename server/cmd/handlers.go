package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/websocket"
)

const (
	// distance radii are in km
	pickupRadius = 5
	destRadius   = 12
)

var (
	activeCarpools = map[string]string{}
)

func PostSignin(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	// decode body and check if json has required data
	var credentials SigninBody
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]any{"ok": false, "user": nil})
		return
	} else if credentials.Email == "" || credentials.Password == "" || len([]byte(credentials.Password)) > 72 {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]any{"ok": false, "user": nil})
		return
	}

	// queries for user in database
	var (
		username, firstname, lastname, email string
		salary                               int64
		passwordHash                         []byte
	)
	row := db.QueryRowContext(
		r.Context(),
		`SELECT username, firstname, lastname, email, salary, passwordHash 
		FROM users 
		WHERE email = ?`,
		credentials.Email,
	)
	err := row.Scan(&username, &firstname, &lastname, &email, &salary, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			enc.Encode(map[string]any{"ok": false, "user": nil})
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false, "user": nil})
		return
	}

	// compare password and hash
	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(credentials.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		enc.Encode(map[string]any{"ok": false, "user": nil})
		return
	}

	sessionManager.Put(r.Context(), "username", username)
	enc.Encode(map[string]any{"ok": false, "user": UserInfo{username, firstname, lastname, email, salary}})
}

func PostSignup(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	// fetch signup body and check data format
	var userInfo SignupBody
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]any{"ok": false, "user": nil})
		return
	}
	if ok := userInfo.CheckFilled(); !ok {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]any{"ok": false, "user": nil})
		return
	}

	// if doesn't exist create new user and send cookie
	pwHash, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false, "user": nil})
		return
	}

	_, err = db.ExecContext(
		r.Context(),
		"INSERT INTO users VALUES (?,?,?,?,?,?,?)",
		userInfo.Username,
		userInfo.Firstname,
		userInfo.LastName,
		userInfo.Email,
		pwHash,
		userInfo.Salary,
		0,
	)
	var sqliteErr sqlite3.Error
	if err != nil && errors.As(err, &sqliteErr) {
		if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			w.WriteHeader(http.StatusConflict)
			enc.Encode(map[string]any{"ok": false, "user": nil})
			return
		}

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false})
		return
	}

	enc.Encode(map[string]any{
		"ok": true,
		"user": UserInfo{
			userInfo.Username,
			userInfo.Firstname,
			userInfo.LastName,
			userInfo.Email,
			userInfo.Salary,
		}})
	return
}

func PostUpdateSalary(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	username := sessionManager.GetString(r.Context(), "username")
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		enc.Encode(map[string]any{"ok": false})
		return
	}

	var u UpdateSalaryBody
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]any{"ok": false})
		return
	}

	_, err := db.ExecContext(r.Context(), "UPDATE users SET salary = ? WHERE username = ?", u.NewSalary, username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false})
		return
	}

	j := map[string]any{"ok": false}
	enc.Encode(j)
	return
}

func PostBookRide(w http.ResponseWriter, r *http.Request) {
	// enc := json.NewEncoder(w)

	//	get all bookings that are up to date and have similar

	// check all carpools from bookings to check if it isn't full

	// create a booking record

}

func WsListenCarpool(ws *websocket.Conn) {
	var cred ListenCarpoolCredentials
	err := json.NewDecoder(ws).Decode(&cred)
	if err != nil {
		log.Printf("error reading ws credentials: %v", err)
	}
	if !cred.CheckFilled() {
		log.Println("empty credentials")
	}

}

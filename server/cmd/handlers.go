package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jftuga/geodist"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const (
	// distance radii are in km
	pickupRadius = 5
	destRadius   = 12

	maxCarpoolSize = 4

	bookingTimeout = 40 * time.Second
	// bookingTimeout = 10 * time.Minute
	carpoolTimeout = bookingTimeout
)

// var (
// 	activeCarpools = map[string]string{}
// )

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
	enc.Encode(map[string]any{"ok": true, "user": UserInfo{username, firstname, lastname, email, salary}})
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

	// insert new user row into users table
	_, err = db.ExecContext(
		r.Context(),
		"INSERT INTO users (username, firstname, lastname, email, passwordhash, salary, accountlevel) VALUES (?,?,?,?,?,?,?)",
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

	// put session cookie and return response
	sessionManager.Put(r.Context(), "username", userInfo.Username)
	enc.Encode(map[string]any{
		"ok": true,
		"user": UserInfo{
			userInfo.Username,
			userInfo.Firstname,
			userInfo.LastName,
			userInfo.Email,
			userInfo.Salary,
		}})
}

func PostUpdateSalary(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	// check session cookie to see if user is authorized
	username := sessionManager.GetString(r.Context(), "username")
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		enc.Encode(map[string]any{"ok": false})
		return
	}

	// decodes request body to get new updated salary
	var u UpdateSalaryBody
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]any{"ok": false})
		return
	}

	// update new salary
	_, err := db.ExecContext(r.Context(), "UPDATE users SET salary = ? WHERE username = ?", u.NewSalary, username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false})
		return
	}

	// return response
	enc.Encode(map[string]any{"ok": false})
}

func PostBookRide(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)

	// check session cookie to see if user is authorized
	username := sessionManager.GetString(r.Context(), "username")
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		enc.Encode(map[string]any{"ok": false, "cost": nil})
		return
	}

	var brBody BookRideBody
	if err := json.NewDecoder(r.Body).Decode(&brBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]any{"ok": false, "cost": nil})
		return
	}

	// get user id
	var reqUserid, salary int64
	row := db.QueryRowContext(r.Context(), "SELECT id, salary FROM users WHERE username = ?", username)
	err := row.Scan(&reqUserid, &salary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]any{"ok": false, "cost": nil})
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false, "cost": nil})
		return
	}

	//	get all bookings that are up to date and have similar distances(10min = 600s)
	curtime := time.Now()

	bookingsNow := []Booking{}
	var (
		rowid, userid, carpoolid, createdAt, expiresAt int64
		pickuplat, pickuplon, destlat, destlon, cost   float64
	)
	rows, err := db.QueryContext(r.Context(), "SELECT rowid, * FROM bookings WHERE ? < expiresAt", curtime.Unix())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false, "cost": nil})
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(
			&rowid,
			&userid,
			&carpoolid,
			&pickuplat,
			&pickuplon,
			&destlat,
			&destlon,
			&createdAt,
			&expiresAt,
			&cost,
		); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(map[string]any{"ok": false, "cost": nil})
			return
		}

		if userid == reqUserid {
			w.WriteHeader(http.StatusConflict)
			enc.Encode(map[string]any{"ok": false, "cost": nil})
			return
		}

		bookingsNow = append(bookingsNow, Booking{
			rowid,
			userid,
			carpoolid,
			pickuplat,
			pickuplon,
			destlat,
			destlon,
			createdAt,
			expiresAt,
		})

	}

	nearBookings := []Booking{}
	curPickupPos := geodist.Coord{Lat: brBody.PickupLat, Lon: brBody.PickupLon}
	curDestPos := geodist.Coord{Lat: brBody.DestLat, Lon: brBody.DestLon}
	for _, b := range bookingsNow {
		nPickupPos := geodist.Coord{Lat: b.PickupLat, Lon: b.PickupLon}
		_, pickUpDist, err := geodist.VincentyDistance(curPickupPos, nPickupPos)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(map[string]any{"ok": false, "cost": nil})
			return
		}

		nDestPos := geodist.Coord{Lat: b.DestLat, Lon: b.DestLon}
		_, destDist, err := geodist.VincentyDistance(curDestPos, nDestPos)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(map[string]any{"ok": false, "cost": nil})
			return
		}

		if pickUpDist <= pickupRadius && destDist <= destRadius {
			nearBookings = append(nearBookings, b)
		}

	}

	// check all carpools from bookings to check if it isn't full
	carpoolIds := map[int64]struct{}{}
	for _, b := range nearBookings {
		carpoolIds[b.CarpoolId] = struct{}{}
	}

	var chosenCarpoolId int64 = -1
	var size int64
	for id := range carpoolIds {
		row = db.QueryRowContext(
			r.Context(),
			"SELECT size FROM carpools WHERE id = ?",
			id,
		)
		if err := row.Scan(&size); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(map[string]any{"ok": false, "cost": nil})
			return
		}

		if size < maxCarpoolSize {
			chosenCarpoolId = id
			break
		}
	}

	if chosenCarpoolId == -1 {
		res, err := db.ExecContext(
			r.Context(),
			"INSERT INTO carpools (createdAt, expiresAt, size) VALUES (?,?,?)",
			curtime.Unix(),
			curtime.Add(carpoolTimeout).Unix(),
			1,
		)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(map[string]any{"ok": false, "cost": nil})
			return
		}

		x, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(map[string]any{"ok": false, "cost": nil})
			return
		}
		chosenCarpoolId = x
	}

	_, dist, err := geodist.VincentyDistance(curPickupPos, curDestPos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false, "cost": nil})
		return
	}
	cost = calcCost(dist, salary)

	// create a booking record
	_, err = db.ExecContext(
		r.Context(),
		`INSERT INTO bookings (
			userid,
			carpoolid,
			pickuplat, 
			pickuplon, 
			destlat, 
			destlon, 
			createdAt, 
			expiresAt,
			cost
		) VALUES (?,?,?,?,?,?,?,?,?)`,
		reqUserid,
		chosenCarpoolId,
		curPickupPos.Lat,
		curPickupPos.Lon,
		curDestPos.Lat,
		curDestPos.Lon,
		curtime.Unix(),
		curtime.Add(carpoolTimeout).Unix(),
		cost,
	)
	if err != nil {
		fmt.Println("hi")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(map[string]any{"ok": false, "cost": nil})
		return
	}

	enc.Encode(map[string]any{"ok": true, "cost": cost})
}

// func WsListenCarpool(ws *websocket.Conn) {
// 	var cred ListenCarpoolCredentials
// 	err := json.NewDecoder(ws).Decode(&cred)
// 	if err != nil {
// 		log.Printf("error reading ws credentials: %v", err)
// 	}
// 	if !cred.CheckFilled() {
// 		log.Println("empty credentials")
// 	}

// }

func calcCost(dist float64, salary int64) float64 {
	return dist*5 + (0.0001 * float64(salary))
}

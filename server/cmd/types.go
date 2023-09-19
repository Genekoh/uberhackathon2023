package main

import "sync"

type UserInfo struct {
	UserName  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Salary    int64  `json:"salary"`
}

type SigninBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type SignupBody struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Salary    int64  `json:"salary"`
}

func (s SignupBody) CheckFilled() bool {
	if s.Username == "" || s.Firstname == "" || s.LastName == "" || s.Email == "" || s.Password == "" || s.Salary == 0 {
		return false
	}
	return true
}

type UpdateSalaryBody struct {
	NewSalary int64 `json:"newSalary"`
}

type BookRideBody struct {
	PickupLat float64 `json:"pickuplat"`
	PickupLon float64 `json:"pickuplon"`
	DestLat   float64 `json:"destlat"`
	DestLon   float64 `json:"destlon"`
}

type Booking struct {
	Id        int64
	UserId    int64
	CarpoolId int64
	PickupLat float64
	PickupLon float64
	DestLat   float64
	DestLon   float64
	CreatedAt int64
	ExpiresAt int64
}

type ListenCarpoolCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Carpool struct {
	UpdateChan chan string
	Users      []string
	Mutex      sync.Mutex
}

func (c *Carpool) AddUser(username string) {
	c.Mutex.Lock()
	c.Users = append(c.Users, username)
	go func() {
		c.UpdateChan <- username
	}()
	c.Mutex.Unlock()
}

func (c ListenCarpoolCredentials) CheckFilled() bool {
	if c.Username == "" || c.Password == "" {
		return false
	}

	return true
}

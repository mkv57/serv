package domain

import "time"

type Session struct {
	Id         int
	UserID     int
	Token      string
	Ip         string
	CreatedAt  time.Time
	User_Agent string
}

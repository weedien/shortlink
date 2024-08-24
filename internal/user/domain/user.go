package domain

import "time"

type User struct {
	id         int
	name       string
	password   string
	realName   string
	email      string
	phone      string
	deleteTime time.Time
}

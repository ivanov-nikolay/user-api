package entity

import "time"

type User struct {
	ID         int64
	Name       string
	Surname    string
	Patronymic string
	Gender     string
	Status     string
	Birthday   time.Time
	JoinDate   time.Time
}

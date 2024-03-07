package dto

import (
	"time"

	"github.com/ivanov-nikolay/user-api/internal/entity"
)

type UserDB struct {
	ID         int64
	Name       string
	Surname    string
	Patronymic *string
	Gender     string
	Status     string
	Birthday   *time.Time
	JoinDate   time.Time
}

func (u *UserDB) ConvertToUser() entity.User {
	var patronymic string
	if u.Patronymic != nil {
		patronymic = *u.Patronymic
	} else {
		patronymic = ""
	}

	var bDay time.Time
	if u.Birthday != nil {
		bDay = *u.Birthday
	} else {
		bDay = time.Time{}
	}
	return entity.User{
		ID:         u.ID,
		Name:       u.Name,
		Surname:    u.Surname,
		Patronymic: patronymic,
		Gender:     u.Gender,
		Status:     u.Status,
		Birthday:   bDay,
		JoinDate:   u.JoinDate,
	}
}

package dto

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ivanov-nikolay/user-api/internal/entity"
)

type UserUpdate struct {
	ID         int64     `json:"id" valid:"required"`
	Name       string    `json:"name" valid:"required,length(2|30),matches(^[A-Z][a-z]+$)"`
	Surname    string    `json:"surname" valid:"required,length(2|30),matches(^[A-Z][a-z]+$)"`
	Patronymic string    `json:"patronymic" valid:"optional,length(2|30),matches(^[A-Z][a-z]+$)"`
	Gender     string    `json:"gender" valid:"required,in(male|female)"`
	Status     string    `json:"status" valid:"required,in(active|banned|deleted)"`
	Birthday   time.Time `json:"b_day" valid:"optional"`
}

func (u *UserUpdate) Validate() []string {
	_, err := govalidator.ValidateStruct(u)
	return collectErrors(err)
}

func (u *UserUpdate) ConvertToUser() entity.User {
	return entity.User{
		ID:         u.ID,
		Name:       u.Name,
		Surname:    u.Surname,
		Patronymic: u.Patronymic,
		Gender:     u.Gender,
		Status:     u.Status,
		Birthday:   u.Birthday,
	}
}

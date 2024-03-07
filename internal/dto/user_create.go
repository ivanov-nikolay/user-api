package dto

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ivanov-nikolay/user-api/internal/entity"
)

type UserCreate struct {
	Name       string    `json:"name" valid:"required,length(2|30),matches(^[A-Z][a-z]+$)"`
	Surname    string    `json:"surname" valid:"required,length(2|30),matches(^[A-Z][a-z]+$)"`
	Patronymic string    `json:"patronymic" valid:"optional,length(2|30),matches(^[A-Z][a-z]+$)"`
	Gender     string    `json:"gender" valid:"required,in(male|female)"`
	Status     string    `json:"status" valid:"required,in(active|banned|deleted)"`
	Birthday   time.Time `json:"b_day" valid:"optional"`
}

func (u *UserCreate) Validate() []string {
	_, err := govalidator.ValidateStruct(u)
	return collectErrors(err)
}

func collectErrors(err error) []string {
	validationErrors := make([]string, 0)
	if err == nil {
		return validationErrors
	}
	var allErrs govalidator.Errors
	if errors.As(err, &allErrs) {
		for _, fld := range allErrs {
			validationErrors = append(validationErrors, fld.Error())
		}
	}
	return validationErrors
}

func (u *UserCreate) ConvertToUser() entity.User {
	return entity.User{
		Name:       u.Name,
		Surname:    u.Surname,
		Patronymic: u.Patronymic,
		Gender:     u.Gender,
		Status:     u.Status,
		Birthday:   u.Birthday,
	}
}

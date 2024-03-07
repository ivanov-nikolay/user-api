package usecase

import (
	"fmt"
	"time"

	"github.com/ivanov-nikolay/user-api/internal/entity"
	"github.com/ivanov-nikolay/user-api/internal/filters"
	"github.com/ivanov-nikolay/user-api/internal/storage"
)

type UserUseCase interface {
	CreateUserUseCase(user entity.User) (*entity.User, error)
	DeleteUserUseCase(ID int64) (bool, error)
	UpdateUserUseCase(user entity.User) (*entity.User, error)
	GetUserByIDUseCase(ID int64) (*entity.User, error)
	SearchUsersUseCase(filters filters.Filter) ([]entity.User, error)
}

type AppUseCase struct {
	s storage.Storage
}

func New(s storage.Storage) *AppUseCase {
	return &AppUseCase{s: s}
}

func (au *AppUseCase) CreateUserUseCase(user entity.User) (*entity.User, error) {
	user.JoinDate = time.Now()
	ID, err := au.s.CreateUserStorage(user)
	if err != nil {
		return nil, fmt.Errorf("storage error: %s", err)
	}
	user.ID = ID
	return &user, nil
}

func (au *AppUseCase) DeleteUserUseCase(ID int64) (bool, error) {
	isDeleted, err := au.s.DeleteUserStorage(ID)
	if err != nil {
		return false, fmt.Errorf("storage error: %s", err)
	}
	return isDeleted, nil
}

func (au *AppUseCase) UpdateUserUseCase(user entity.User) (*entity.User, error) {
	wasUpdated, err := au.s.UpdateUserStorage(user)
	if err != nil {
		return nil, fmt.Errorf("storage error: %s", err)
	}
	if !wasUpdated {
		return nil, nil
	}
	return &user, nil
}

func (au *AppUseCase) GetUserByIDUseCase(ID int64) (*entity.User, error) {
	user, err := au.s.GetUserByIDStorage(ID)
	if err != nil {
		return nil, fmt.Errorf("storage error: %s", err)
	}
	return user, nil

}

func (au *AppUseCase) SearchUsersUseCase(filter filters.Filter) ([]entity.User, error) {
	users, err := au.s.SearchUsersStorage(filter)
	if err != nil {
		return nil, fmt.Errorf("storage error: %s", err)
	}
	return users, nil

}

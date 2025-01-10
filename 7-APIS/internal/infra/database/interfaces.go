package database

import "github.com/diogokimisima/goexpert/7-APIS/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

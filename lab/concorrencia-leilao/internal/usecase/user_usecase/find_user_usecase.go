package user_usecase

import (
	"context"

	"github.com/diogokimisima/fullcycle-auction/internal/entity/user_entity"
	"github.com/diogokimisima/fullcycle-auction/internal/internal_error"
)

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type UserUseCaseInterface interface {
	FindUserById(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError)
}

func (u *UserUseCase) FindUserById(
	ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError) {

	userEntity, err := u.UserRepository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:       userEntity.Id,
		Username: userEntity.Name,
	}, nil
}

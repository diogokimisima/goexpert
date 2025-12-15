package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/diogokimisima/fullcycle-auction/configuration/logger"
	"github.com/diogokimisima/fullcycle-auction/internal/entity/user_entity"
	"github.com/diogokimisima/fullcycle-auction/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (ur *UserRepository) FindUserById(
	ctx context.Context, userId string) (*user_entity.User, error) {
	filter := bson.M{"_id": userId}

	var userEntityMongo UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			errMsg := fmt.Sprintf("user with id %s not found", userId)
			logger.Error(errMsg, mongo.ErrNoDocuments)
			return nil, internal_error.NewNotFoundError(errMsg)
		}

		logger.Error(fmt.Sprintf("Error trying to find user by userId %s: %v", userId, err), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error trying to find user by userId %s: %v", userId, err))
	}

	userEntity := &user_entity.User{
		Id:   userEntityMongo.Id,
		Name: userEntityMongo.Name,
	}

	return userEntity, nil
}

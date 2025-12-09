package user

import (
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

// func (ur *UserRepository) FindUserById(
// 	ctx context.Context, userId string) (*user_entity.User, error) {
// 	filter := bson.M{"_id": userId}

// }

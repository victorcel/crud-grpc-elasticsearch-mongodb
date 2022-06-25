package repository

import (
	"context"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) (string, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	UpdateUser(ctx context.Context, id string, name string) (*mongo.UpdateResult, error)
	DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error)
	Close() error
}

var implementationUser UserRepository

func SetUserRepository(repository UserRepository) {
	implementationUser = repository
}

func InsertUser(ctx context.Context, user *models.User) (string, error) {
	return implementationUser.InsertUser(ctx, user)
}

func GetUserByID(ctx context.Context, id string) (models.User, error) {
	return implementationUser.GetUserByID(ctx, id)
}

func UpdateUser(ctx context.Context, id string, name string) (*mongo.UpdateResult, error) {
	return implementationUser.UpdateUser(ctx, id, name)
}

func DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	return implementationUser.DeleteUser(ctx, id)
}

func Close() error {
	return implementationUser.Close()
}

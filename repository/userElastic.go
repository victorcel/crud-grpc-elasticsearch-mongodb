package repository

import (
	"context"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/models"
)

type UserElasticRepository interface {
	InsertUserElastic(ctx context.Context, user *models.UserElastic) (string, error)
	GetUserElasticByID(ctx context.Context, id string) (*models.UserElastic, error)
	UpdateUserElastic(ctx context.Context, user models.UserElastic) error
	DeleteUserElastic(ctx context.Context, id string) error
}

var implementationUserElastic UserElasticRepository

func SetUserElasticRepository(repository UserElasticRepository) {
	implementationUserElastic = repository
}

func InsertUserElastic(ctx context.Context, user *models.UserElastic) (string, error) {
	return implementationUserElastic.InsertUserElastic(ctx, user)
}

func GetUserElasticByID(ctx context.Context, id string) (*models.UserElastic, error) {
	return implementationUserElastic.GetUserElasticByID(ctx, id)
}

func UpdateUserElastic(ctx context.Context, user models.UserElastic) error {
	return implementationUserElastic.UpdateUserElastic(ctx, user)
}

func DeleteUserElastic(ctx context.Context, id string) error {
	return implementationUserElastic.DeleteUserElastic(ctx, id)
}

package database

import (
	"context"
	"fmt"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	CollectionName = "users"
)

type MongoDbRepository struct {
	db *mongo.Database
}

func NewMongoDbRepository(url string) (*MongoDbRepository, error) {

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error connect mongodb", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error ping mongodb", err)
	}

	db := client.Database("grpc-vozy")
	return &MongoDbRepository{db}, nil
}

func (repository *MongoDbRepository) Close() error {
	return repository.db.Client().Disconnect(context.Background())
}

func (repository *MongoDbRepository) InsertUser(ctx context.Context, user *models.User) (string, error) {
	result, err := repository.db.Collection(CollectionName).InsertOne(ctx, user)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result.InsertedID.(primitive.ObjectID).Hex()), err
}

func (repository *MongoDbRepository) GetUserByID(ctx context.Context, id string) (models.User, error) {

	var user models.User

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return user, err
	}

	err = repository.db.Collection(CollectionName).FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repository *MongoDbRepository) UpdateUser(ctx context.Context, id string, name string) (
	*mongo.UpdateResult, error,
) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{{"name", name}}}}
	updateOne, err := repository.db.Collection(CollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return updateOne, nil
}

func (repository *MongoDbRepository) DeleteUser(ctx context.Context, id string) (
	*mongo.DeleteResult, error,
) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	deleteOne, err := repository.db.Collection(CollectionName).DeleteOne(ctx, bson.D{{"_id", objectId}})
	if err != nil {
		return nil, err
	}
	return deleteOne, nil
}

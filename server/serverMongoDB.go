package server

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/database"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/models"
	userProtoBuff "github.com/victorcel/crud-grpc-elasticsearch-mongodb/proto"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

type ServerMongoDB struct {
	repository repository.UserRepository
	userProtoBuff.UnimplementedUserServiceServer
	validate *validator.Validate
}

func NewServerMongoDB(repository repository.UserRepository) *ServerMongoDB {
	validate := validator.New()
	return &ServerMongoDB{
		repository: repository,
		validate:   validate,
	}
}

func InitializeMongoDBServer() {
	fmt.Println("Iniciado servidor MongoDB GPRC...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	listenServer, err := net.Listen("tcp", ":"+os.Getenv("PORT_MONGODB"))
	if err != nil {
		log.Fatal(err)
	}

	repository, err := database.NewMongoDbRepository(os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	serverMain := NewServerMongoDB(repository)

	serverGrpc := grpc.NewServer()
	userProtoBuff.RegisterUserServiceServer(serverGrpc, serverMain)

	reflection.Register(serverGrpc)

	fmt.Println("ServerMongoDb run: " + os.Getenv("PORT_MONGODB"))

	if err := serverGrpc.Serve(listenServer); err != nil {
		log.Fatalf("Error serving: %s", err.Error())
	}
}

func (s *ServerMongoDB) InsertUser(ctx context.Context, request *userProtoBuff.User) (*userProtoBuff.UserByIdResponse, error) {

	user := &models.User{
		ID:    primitive.NewObjectID(),
		Name:  request.GetName(),
		Email: request.GetEmail(),
		Ega:   int(request.GetEga()),
	}

	err := s.validate.Struct(user)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userResponse, err := s.repository.InsertUser(ctx, user)

	if err != nil {
		return nil, err
	}

	return &userProtoBuff.UserByIdResponse{
		Id: userResponse,
	}, nil
}

func (s *ServerMongoDB) GetUserByID(ctx context.Context, request *userProtoBuff.UserRequest) (*userProtoBuff.User, error) {

	userResponse, err := s.repository.GetUserByID(ctx, request.GetId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &userProtoBuff.User{
		Id:    userResponse.ID.Hex(),
		Name:  userResponse.Name,
		Email: userResponse.Email,
		Ega:   int64(userResponse.Ega),
	}, nil
}

func (s *ServerMongoDB) UpdateUser(ctx context.Context, request *userProtoBuff.User) (*userProtoBuff.UserResponse, error) {
	_, err := s.repository.UpdateUser(ctx, request.GetId(), request.GetName())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &userProtoBuff.UserResponse{Result: true}, nil
}

func (s *ServerMongoDB) DeleteUser(ctx context.Context, request *userProtoBuff.UserRequest) (
	*userProtoBuff.UserResponse, error,
) {
	_, err := s.repository.DeleteUser(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &userProtoBuff.UserResponse{Result: true}, nil
}

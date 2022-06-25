package server

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/database"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/models"
	userProtoBuff "github.com/victorcel/crud-grpc-elasticsearch-mongodb/proto"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/repository"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/search"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

type ServerElasticDB struct {
	userProtoBuff.UnimplementedUserServiceServer
	validate *validator.Validate
}

func NewServerElasticDB() *ServerElasticDB {
	validate := validator.New()
	return &ServerElasticDB{
		validate: validate,
	}
}

func InitializeElasticServer() {
	fmt.Println("Iniciado servidor Elastic...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	listenServer, err := net.Listen("tcp", ":"+os.Getenv("PORT_ELASTICSEARCH"))
	if err != nil {
		log.Fatal(err)
	}
	elastic, err := search.New([]string{os.Getenv("ELASTICSEARCH_ADDRESS")})
	if err != nil {
		log.Fatalln(err)
	}
	if err := elastic.CreateIndex("user"); err != nil {
		log.Fatalln(err)
	}

	storage, err := database.NewUserStorage(*elastic)
	if err != nil {
		log.Fatalln(err)
	}

	serverMain := NewServerElasticDB()

	repository.SetUserElasticRepository(&storage)

	serverGrpc := grpc.NewServer()

	userProtoBuff.RegisterUserServiceServer(serverGrpc, serverMain)

	reflection.Register(serverGrpc)

	fmt.Println("ServerMongoDb run: " + os.Getenv("PORT_ELASTICSEARCH"))

	if err := serverGrpc.Serve(listenServer); err != nil {
		log.Fatalf("Error serving: %s", err.Error())
	}
}

func (s *ServerElasticDB) InsertUser(ctx context.Context, request *userProtoBuff.User) (*userProtoBuff.UserByIdResponse, error) {

	user := &models.UserElastic{
		ID:    primitive.NewObjectID().Hex(),
		Name:  request.GetName(),
		Email: request.GetEmail(),
		Ega:   int(request.GetEga()),
	}

	err := s.validate.Struct(user)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userResponse, err := repository.InsertUserElastic(ctx, user)

	if err != nil {
		return nil, err
	}

	return &userProtoBuff.UserByIdResponse{
		Id: userResponse,
	}, nil
}

func (s *ServerElasticDB) GetUserByID(ctx context.Context, request *userProtoBuff.UserRequest) (*userProtoBuff.User, error) {

	userResponse, err := repository.GetUserElasticByID(ctx, request.GetId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &userProtoBuff.User{
		Id:    userResponse.ID,
		Name:  userResponse.Name,
		Email: userResponse.Email,
		Ega:   int64(userResponse.Ega),
	}, nil
}

func (s *ServerElasticDB) UpdateUser(ctx context.Context, request *userProtoBuff.User) (*userProtoBuff.UserResponse, error) {
	user := &models.UserElastic{
		ID:    request.GetId(),
		Name:  request.GetName(),
		Email: request.GetEmail(),
		Ega:   int(request.GetEga()),
	}

	err := s.validate.Struct(user)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = repository.UpdateUserElastic(ctx, *user)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &userProtoBuff.UserResponse{Result: true}, nil
}

func (s *ServerElasticDB) DeleteUser(ctx context.Context, request *userProtoBuff.UserRequest) (
	*userProtoBuff.UserResponse, error,
) {
	err := repository.DeleteUserElastic(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &userProtoBuff.UserResponse{Result: true}, nil
}

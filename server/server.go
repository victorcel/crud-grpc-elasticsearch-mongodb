package server

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/models"
	userProtoBuff "github.com/victorcel/crud-grpc-elasticsearch-mongodb/proto"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	repository repository.UserRepository
	userProtoBuff.UnimplementedUserServiceServer
	validate *validator.Validate
}

func NewServer(repository repository.UserRepository) *Server {
	validate := validator.New()
	return &Server{
		repository: repository,
		validate:   validate,
	}
}

func (s *Server) InsertUser(ctx context.Context, request *userProtoBuff.User) (*userProtoBuff.UserByIdResponse, error) {

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

func (s *Server) GetUserByID(ctx context.Context, request *userProtoBuff.UserRequest) (*userProtoBuff.User, error) {

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

func (s *Server) UpdateUser(ctx context.Context, request *userProtoBuff.User) (*userProtoBuff.UserResponse, error) {
	_, err := s.repository.UpdateUser(ctx, request.GetId(), request.GetName())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &userProtoBuff.UserResponse{Result: true}, nil
}

func (s *Server) DeleteUser(ctx context.Context, request *userProtoBuff.UserRequest) (
	*userProtoBuff.UserResponse, error,
) {
	_, err := s.repository.DeleteUser(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &userProtoBuff.UserResponse{Result: true}, nil
}

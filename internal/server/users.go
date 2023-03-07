package server

import (
	"context"
	"errors"
	"github.com/practice-sem-2/user-service/internal/models"
	"github.com/practice-sem-2/user-service/internal/pb"
	storage "github.com/practice-sem-2/user-service/internal/storages"
	usecase "github.com/practice-sem-2/user-service/internal/usecases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServer
	ucase *usecase.UseCase
}

var (
	ErrUserNotFound       = status.Error(codes.NotFound, "user with provided username does not exist")
	ErrUserAlreadyExists  = status.Error(codes.AlreadyExists, "user already exists")
	ErrEmailAlreadyExists = status.Error(codes.AlreadyExists, "provided email is already taken")
)

func wrapError(err error) error {
	errorMapper := []struct {
		from error
		to   error
	}{
		{from: storage.ErrUserNotFound, to: ErrUserNotFound},
		{from: storage.ErrUserAlreadyExists, to: ErrUserAlreadyExists},
		{from: storage.ErrEmailAlreadyExists, to: ErrEmailAlreadyExists},
	}

	if err == nil {
		return nil
	}

	for _, mapping := range errorMapper {
		if errors.Is(err, mapping.from) {
			return mapping.to
		}
	}
	return status.Error(codes.Internal, err.Error())
}

func (u *UserServer) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userCreate, err := ParseCreateRequest(r)

	if err != nil {
		return nil, err
	}

	user, err := u.ucase.Users.Create(ctx, &userCreate)

	if err != nil {
		return nil, wrapError(err)
	}

	return &pb.CreateUserResponse{
		User: ToUserData(user),
	}, nil
}

func (u *UserServer) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user *models.User
	var err error

	if r.Username == nil && r.Email == nil {
		return nil, status.Error(codes.InvalidArgument, "Either username or email must be provided")
	} else if r.Username != nil {
		user, err = u.ucase.Users.GetByUsername(ctx, *r.Username)
	} else if r.Email != nil {
		user, err = u.ucase.Users.GetByEmail(ctx, *r.Email)
	}

	if err != nil {
		return nil, wrapError(err)
	}

	return &pb.GetUserResponse{
		User: ToUserData(user),
	}, nil

}

func (u *UserServer) UpdateUser(ctx context.Context, r *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	update := models.UpdateFields{
		Password:  nil,
		Email:     nil,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		AvatarID:  r.AvatarId,
	}

	user, err := u.ucase.Users.Update(ctx, r.Username, update)

	if err != nil {
		return nil, wrapError(err)
	}

	return &pb.UpdateUserResponse{
		User: ToUserData(user),
	}, nil
}

func (u *UserServer) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := u.ucase.Users.Delete(ctx, r.Username)
	return &pb.DeleteUserResponse{}, wrapError(err)
}

func NewUserServer(ucase *usecase.UseCase) *UserServer {
	return &UserServer{
		ucase: ucase,
	}
}

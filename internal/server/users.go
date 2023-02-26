package server

import (
	"context"
	"github.com/practice-sem-2/user-service/internal/pb"
)

type UserServer struct {
	pb.UnimplementedUserServer
}

func (u *UserServer) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserServer) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserServer) UpdateUser(ctx context.Context, t *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserServer) DeleteUser(ctx context.Context, t *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserServer() *UserServer {
	return &UserServer{}
}

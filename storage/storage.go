package storage

import (
	pb "auth-service/genproto/user"
	"auth-service/models"
	"context"
)

type IStorage interface {
	User() IUserStorage
	Close()
}

type IUserStorage interface {
	Add(ctx context.Context, u *models.RegisterRequest) (*models.RegisterResponse, error)
	Read(ctx context.Context, u *pb.ID) (*pb.Profile, error)
	Update(ctx context.Context, u *pb.NewData) (*pb.UpdateResp, error)
	GetDetails(ctx context.Context, email string) (*models.UserDetails, error)
	GetRole(ctx context.Context, id string) (string, error)
}

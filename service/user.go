package service

import (
	pb "auth-service/genproto/user"
	"auth-service/pkg/logger"
	"auth-service/storage"
	"context"
	"log/slog"

	"github.com/pkg/errors"
)

type UserService struct {
	pb.UnimplementedUserServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewUserService(s storage.IStorage) *UserService {
	return &UserService{
		storage: s,
		logger:  logger.NewLogger(),
	}
}

func (s *UserService) GetProfile(ctx context.Context, req *pb.ID) (*pb.Profile, error) {
	s.logger.Info("GetProfile is invoked", slog.Any("request", req))

	resp, err := s.storage.User().Read(ctx, req)
	if err != nil {
		er := errors.Wrap(err, "failed to get profile")
		s.logger.Error(er.Error())
		return nil, er
	}

	s.logger.Info("GetProfile is completed", slog.Any("response", resp))
	return resp, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.NewData) (*pb.UpdateResp, error) {
	s.logger.Info("UpdateProfile is invoked", slog.Any("request", req))

	resp, err := s.storage.User().Update(ctx, req)
	if err != nil {
		er := errors.Wrap(err, "failed to update profile")
		s.logger.Error(er.Error())
		return nil, er
	}

	s.logger.Info("UpdateProfile is completed", slog.Any("response", resp))
	return resp, nil
}

func (s *UserService) ValidateUser(ctx context.Context, req *pb.ID) (*pb.Void, error) {
	s.logger.Info("ValidateUser is invoked", slog.Any("request", req))

	_, err := s.storage.User().Read(ctx, req)
	if err != nil {
		er := errors.Wrap(err, "failed to validate user")
		s.logger.Error(er.Error())
		return nil, er
	}

	s.logger.Info("ValidateUser is completed")
	return &pb.Void{}, nil
}

package handlers

import (
	"context"
	"library-api-user/pkg/token"
	pb "library-api-user/proto/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	payload, err := token.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Success: false,
			AuthId:  0,
			Role:    "",
		}, nil
	}

	return &pb.ValidateResponse{
		Success: true,
		AuthId:  uint64(payload.AuthId),
		Role:    payload.Role,
	}, nil
}

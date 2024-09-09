package authgrpc

import (
	"context"
	"sso/protos/gen/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(
		ctx context.Context,
		userID int64,
	) (bool, error)
}

type ServerAPI struct {
	proto.UnimplementedAuthServer      // просто заглушка для запуска, даже если что-то не реализовано
	auth                          Auth //interface ALWAYS = 2 pointers
}

func Register(gRPC *grpc.Server, auth Auth) {
	proto.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.RegisterResponse{UserId: userID}, nil
}

// Login logs in a user and returns an auth token.
func (s *ServerAPI) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		//TODO: not safe to send internal errors to client
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.LoginResponse{Token: token}, nil
}

// IsAdmin checks whether a user is an admin.
func (s *ServerAPI) IsAdmin(ctx context.Context, req *proto.IsAdminRequest) (*proto.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func (s *ServerAPI) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	panic("implement me")
}

func validateLogin(req *proto.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func validateRegister(req *proto.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateIsAdmin(req *proto.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}

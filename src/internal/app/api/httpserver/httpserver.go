package httpserver

import (
	"context"
	"log/slog"
	userService "service-template/internal/app/cases/user"
	"service-template/pkg/api"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
)

type HttpServer struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *HttpServer {
	return &HttpServer{
		logger: logger,
	}
}

func (s *HttpServer) GetV1Users(ctx context.Context, request api.GetV1UsersRequestObject) (api.GetV1UsersResponseObject, error) {
	service, err := userService.Build(s.logger)
	if err != nil {
		return nil, err
	}
	return service.GetUsers(ctx, request)
}

func (s *HttpServer) GetV2Users(ctx context.Context, request api.GetV2UsersRequestObject) (api.GetV2UsersResponseObject, error) {
	service, err := userService.Build(s.logger)
	if err != nil {
		return nil, err
	}
	return service.GetUsersV2(ctx, request)
}


func (s *HttpServer) PostV1Users(ctx context.Context, request api.PostV1UsersRequestObject) (api.PostV1UsersResponseObject, error) {
	service, err := userService.Build(s.logger)
	if err != nil {
		return nil, err
	}
	return service.CreateUser(ctx, request)
}

func (s *HttpServer) DeleteV1UsersID(ctx context.Context, request api.DeleteV1UsersIDRequestObject) (api.DeleteV1UsersIDResponseObject, error) {
	service, err := userService.Build(s.logger)
	if err != nil {
		return nil, err
	}
	return service.DeleteUser(ctx, request)
}

func (s *HttpServer) GetV1UsersID(ctx context.Context, request api.GetV1UsersIDRequestObject) (api.GetV1UsersIDResponseObject, error) {
	service, err := userService.Build(s.logger)
	if err != nil {
		return nil, err
	}
	return service.GetUser(ctx, request)
}

func (s *HttpServer) PutV1UsersID(ctx context.Context, request api.PutV1UsersIDRequestObject) (api.PutV1UsersIDResponseObject, error) {
	service, err := userService.Build(s.logger)
	if err != nil {
		return nil, err
	}
	return service.UpdateUser(ctx, request)
}


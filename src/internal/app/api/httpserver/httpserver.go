package httpserver

import (
	"context"
	userService "service-template/internal/app/core/services/user"
	"service-template/pkg/api"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
)

type HttpServer struct {
}

func (HttpServer) GetV1Users(ctx context.Context, request api.GetV1UsersRequestObject) (api.GetV1UsersResponseObject, error) {
	service, err := userService.Build()
	if err != nil {
		return nil, err
	}
	return service.GetUsers(ctx, request)
}

func (HttpServer) GetV2Users(ctx context.Context, request api.GetV2UsersRequestObject) (api.GetV2UsersResponseObject, error) {
	service, err := userService.Build()
	if err != nil {
		return nil, err
	}
	return service.GetUsersV2(ctx, request)
}


func (HttpServer) PostV1Users(ctx context.Context, request api.PostV1UsersRequestObject) (api.PostV1UsersResponseObject, error) {
	service, err := userService.Build()
	if err != nil {
		return nil, err
	}
	return service.CreateUser(ctx, request)
}

func (HttpServer) DeleteV1UsersID(ctx context.Context, request api.DeleteV1UsersIDRequestObject) (api.DeleteV1UsersIDResponseObject, error) {
	service, err := userService.Build()
	if err != nil {
		return nil, err
	}
	return service.DeleteUser(ctx, request)
}

func (HttpServer) GetV1UsersID(ctx context.Context, request api.GetV1UsersIDRequestObject) (api.GetV1UsersIDResponseObject, error) {
	service, err := userService.Build()
	if err != nil {
		return nil, err
	}
	return service.GetUser(ctx, request)
}

func (HttpServer) PutV1UsersID(ctx context.Context, request api.PutV1UsersIDRequestObject) (api.PutV1UsersIDResponseObject, error) {
	service, err := userService.Build()
	if err != nil {
		return nil, err
	}
	return service.UpdateUser(ctx, request)
}


package httpserver

import (
	"context"
	"service-template/pkg/api"
)

type HttpServer struct {
}

func (HttpServer) GetV1Users(ctx context.Context, request api.GetV1UsersRequestObject) (api.GetV1UsersResponseObject, error) {
	return api.GetV1Users200JSONResponse{}, nil
}

func (HttpServer) PostV1Users(ctx context.Context, request api.PostV1UsersRequestObject) (api.PostV1UsersResponseObject, error) {
	return api.PostV1Users201JSONResponse{}, nil
}

func (HttpServer) DeleteV1UsersID(ctx context.Context, request api.DeleteV1UsersIDRequestObject) (api.DeleteV1UsersIDResponseObject, error) {
	return api.DeleteV1UsersID204Response{}, nil
}

func (HttpServer) GetV1UsersID(ctx context.Context, request api.GetV1UsersIDRequestObject) (api.GetV1UsersIDResponseObject, error) {
	return api.GetV1UsersID200JSONResponse{}, nil
}

func (HttpServer) PutV1UsersID(ctx context.Context, request api.PutV1UsersIDRequestObject) (api.PutV1UsersIDResponseObject, error) {
	return api.PutV1UsersID200JSONResponse{}, nil
}

func (HttpServer) GetV2Users(ctx context.Context, request api.GetV2UsersRequestObject) (api.GetV2UsersResponseObject, error) {
	return api.GetV2Users200JSONResponse{}, nil
}

package user

import (
	"context"
	"math"
	"service-template/internal/app/core/domain/common"
	domainUser "service-template/internal/app/core/domain/user"
	"service-template/internal/config"
	"service-template/pkg/api"
)

type UserService struct {
	// UserRepo *userRepository.UserRepository
	UserRepo domainUser.IRepository
	Cfg      config.Config
}

func NewUserService(userRepo domainUser.IRepository, cfg config.Config) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (service *UserService) GetUsers(ctx context.Context, request api.GetV1UsersRequestObject) (api.GetV1UsersResponseObject, error) {
	users, total, err := service.UserRepo.GetUsers(ctx, request.Params.Limit, request.Params.Page)
	if err != nil {
		return nil, err
	}

	var apiUsers []api.User
	for _, v := range users {
		apiUser := api.User{
			ID:       v.Id().Value(),
			Username: v.UserName().Value(),
		}
		apiUsers = append(apiUsers, apiUser)
	}

	totalPages := service.calculateTotalPages(total, request.Params.Limit)

	response := api.GetV1Users200JSONResponse{
		Data: apiUsers,
		Pagination: api.PaginationInfo{
			Limit:      request.Params.Limit,
			Page:       request.Params.Page,
			Total:      total,
			TotalPages: totalPages,
		},
	}
	return response, nil
}

func (service *UserService) GetUsersV2(ctx context.Context, request api.GetV2UsersRequestObject) (api.GetV2UsersResponseObject, error) {
	users, total, err := service.UserRepo.GetUsers(ctx, request.Params.Limit, request.Params.Page)
	if err != nil {
		return nil, err
	}

	var apiUsers []api.UserV2
	for _, v := range users {
		apiUser := api.UserV2{
			ID:       v.Id().Value(),
			Username: v.UserName().Value(),
			Email:    v.Email().Value(),
		}
		apiUsers = append(apiUsers, apiUser)
	}

	totalPages := service.calculateTotalPages(total, request.Params.Limit)

	response := api.GetV2Users200JSONResponse{
		Data: apiUsers,
		Pagination: api.PaginationInfo{
			Limit:      request.Params.Limit,
			Page:       request.Params.Page,
			Total:      total,
			TotalPages: totalPages,
		},
	}
	return response, nil
}

func (service *UserService) CreateUser(ctx context.Context, request api.PostV1UsersRequestObject) (api.PostV1UsersResponseObject, error) {
	user, err := domainUser.New(
		request.Body.Username,
		string(request.Body.Email),
	)
	if err != nil {
		return nil, err
	}

	err = service.UserRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	response := api.PostV1Users201JSONResponse{
		Data: &api.User{
			ID:       user.Id().Value(),
			Username: user.UserName().Value(),
		},
	}
	return response, nil
}

func (service *UserService) DeleteUser(ctx context.Context, request api.DeleteV1UsersIDRequestObject) (api.DeleteV1UsersIDResponseObject, error) {
	err := service.UserRepo.Delete(ctx, request.ID)
	if err != nil {
		return nil, err
	}
	return api.DeleteV1UsersID204Response{}, nil
}

func (service *UserService) GetUser(ctx context.Context, request api.GetV1UsersIDRequestObject) (api.GetV1UsersIDResponseObject, error) {
	user, err := service.UserRepo.GetUser(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	response := api.GetV1UsersID200JSONResponse{
		Data: &api.User{
			ID:       user.Id().Value(),
			Username: user.UserName().Value(),
		},
	}
	return response, nil
}

func (service *UserService) UpdateUser(ctx context.Context, request api.PutV1UsersIDRequestObject) (api.PutV1UsersIDResponseObject, error) {
	user, err := service.UserRepo.GetUser(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	userEmail, err := common.NewEmail(string(request.Body.Email))
	if err != nil {
		return nil, err
	}

	user.SetEmail(userEmail)

	err = service.UserRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	response := api.PutV1UsersID200JSONResponse{
		Data: &api.User{
			ID:       user.Id().Value(),
			Username: user.UserName().Value(),
		},
	}
	return response, nil
}

func (*UserService) calculateTotalPages(total int, limit int) int {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return totalPages
}

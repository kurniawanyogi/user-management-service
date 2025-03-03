package service

import (
	"user-management-service/service/user"
)

type IRegistry interface {
	GetUserService() user.IUserService
}

type Registry struct {
	userService user.IUserService
}

func NewRegistry(
	userService user.IUserService,
) *Registry {
	return &Registry{
		userService: userService,
	}
}

func (r *Registry) GetUserService() user.IUserService {
	return r.userService
}

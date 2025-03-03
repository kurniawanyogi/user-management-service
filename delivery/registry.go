package delivery

import (
	"user-management-service/delivery/user"
)

type IRegistry interface {
	GetUserDelivery() user.IUserDelivery
}

type Registry struct {
	userDelivery user.IUserDelivery
}

func NewRegistry(userDelivery user.IUserDelivery) *Registry {
	return &Registry{
		userDelivery: userDelivery,
	}
}

func (r *Registry) GetUserDelivery() user.IUserDelivery {
	return r.userDelivery
}

package repository

import (
	"github.com/jmoiron/sqlx"
	"user-management-service/repository/user"
)

// @Notice: Register your repositories here

type IRegistry interface {
	GetUserRepository() user.IUserRepository
}

type registry struct {
	dbMaster       *sqlx.DB
	userRepository user.IUserRepository
}

func NewRegistry(
	dbMaster *sqlx.DB,
	userRepository user.IUserRepository,
) *registry {
	return &registry{
		dbMaster:       dbMaster,
		userRepository: userRepository,
	}
}

func (r registry) GetUserRepository() user.IUserRepository {
	return r.userRepository
}

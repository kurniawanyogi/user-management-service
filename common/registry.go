package common

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

// IRegistry TODO Implement redis etc here
type IRegistry interface {
	GetValidator() *validator.Validate
}

type registry struct {
	mu        *sync.Mutex
	validator *validator.Validate
}

func WithValidator(validator *validator.Validate) Option {
	return func(s *registry) {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.validator = validator
	}
}

type Option func(r *registry)

func NewRegistry(
	options ...Option,
) IRegistry {
	registry := &registry{mu: &sync.Mutex{}}

	for _, option := range options {
		option(registry)
	}

	return registry
}

func (r *registry) GetValidator() *validator.Validate {
	return r.validator
}

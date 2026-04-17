package user

import (
	"context"
	"fmt"
)

type UserService struct {
	repository Repository
}

func NewService(repository Repository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetUsers(ctx context.Context) ([]User, error) {
	users, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Fetching all users: %w", err)
	}
	return users, nil
}

func (s *UserService) GetUserById(ctx context.Context, id string) (User, error) {
	user, err := s.repository.GetById(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("Fetching user by id: %w", err)
	}
	return user, nil
}

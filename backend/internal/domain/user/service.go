package user

import (
	"context"
	"fmt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	users, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Fetching all users: %w", err)
	}
	return users, nil
}

func (s *Service) GetUserById(ctx context.Context, id string) (User, error) {
	user, err := s.repository.GetById(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("Fetching user by id: %w", err)
	}
	return user, nil
}

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

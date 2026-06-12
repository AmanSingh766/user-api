package service

import (
	"context"
	"time"

	"github.com/dinesh/user-api/internal/models"
	"github.com/dinesh/user-api/internal/repository"
	"github.com/dinesh/user-api/internal/logger"
	"go.uber.org/zap"
)

const dobLayout = "2006-01-02"

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error)
	GetUser(ctx context.Context, id int32) (models.UserWithAgeResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, page, limit int) (models.PaginatedUsersResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CalculateAge returns the age in full years based on dob and today's date.
// It correctly handles the case where the birthday hasn't occurred yet this year.
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	// Birthday hasn't happened yet this year
	if now.Month() < dob.Month() ||
		(now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	return years
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		logger.Log.Error("failed to create user", zap.Error(err))
		return models.UserResponse{}, err
	}

	logger.Log.Info("user created", zap.Int32("id", user.ID), zap.String("name", user.Name))
	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format(dobLayout),
	}, nil
}

func (s *userService) GetUser(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.UserWithAgeResponse{}, err
	}

	logger.Log.Info("fetched user", zap.Int32("id", user.ID))
	return models.UserWithAgeResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format(dobLayout),
		Age:  CalculateAge(user.Dob),
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.DOB)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		logger.Log.Error("failed to update user", zap.Int32("id", id), zap.Error(err))
		return models.UserResponse{}, err
	}

	logger.Log.Info("user updated", zap.Int32("id", user.ID))
	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format(dobLayout),
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	logger.Log.Info("user deleted", zap.Int32("id", id))
	return nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) (models.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := int32((page - 1) * limit)
	users, err := s.repo.List(ctx, int32(limit), offset)
	if err != nil {
		logger.Log.Error("failed to list users", zap.Error(err))
		return models.PaginatedUsersResponse{}, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return models.PaginatedUsersResponse{}, err
	}

	var list []models.UserWithAgeResponse
	for _, u := range users {
		list = append(list, models.UserWithAgeResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.Dob.Format(dobLayout),
			Age:  CalculateAge(u.Dob),
		})
	}
	if list == nil {
		list = []models.UserWithAgeResponse{}
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return models.PaginatedUsersResponse{
		Data:       list,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

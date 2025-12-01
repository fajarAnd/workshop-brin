package services

import (
	"context"
	"log"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"
	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/repositories"

	"github.com/google/uuid"
)

type UserService interface {
	IsUserEligible(ctx context.Context, phone string) (bool, error)
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetEligibleUsers(ctx context.Context) ([]*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) IsUserEligible(ctx context.Context, phone string) (bool, error) {
	log.Printf("[UserService] Checking eligibility for phone: %s", phone)

	eligible, err := s.userRepo.IsEligible(ctx, phone)
	if err != nil {
		log.Printf("[UserService] Failed to check eligibility for phone %s: %v", phone, err)
		return false, err
	}

	log.Printf("[UserService] Phone %s eligibility: %t", phone, eligible)
	return eligible, nil
}

func (s *userService) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	log.Printf("[UserService] Getting user by phone: %s", phone)

	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		log.Printf("[UserService] Failed to get user by phone %s: %v", phone, err)
		return nil, err
	}

	log.Printf("[UserService] Found user: %s (ID: %s)", user.Name, user.ID.String())
	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	log.Printf("[UserService] Getting user by ID: %s", id.String())

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("[UserService] Failed to get user by ID %s: %v", id.String(), err)
		return nil, err
	}

	log.Printf("[UserService] Found user: %s (Phone: %s)", user.Name, user.Phone)
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	log.Printf("[UserService] Creating user: %s (Phone: %s)", req.Name, req.Phone)

	user, err := s.userRepo.Create(ctx, req)
	if err != nil {
		log.Printf("[UserService] Failed to create user: %v", err)
		return nil, err
	}

	log.Printf("[UserService] User created successfully: %s (ID: %s)", user.Name, user.ID.String())
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	log.Printf("[UserService] Updating user ID: %s", id.String())

	user, err := s.userRepo.Update(ctx, id, req)
	if err != nil {
		log.Printf("[UserService] Failed to update user ID %s: %v", id.String(), err)
		return nil, err
	}

	log.Printf("[UserService] User updated successfully: %s (ID: %s)", user.Name, user.ID.String())
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	log.Printf("[UserService] Deleting user ID: %s", id.String())

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		log.Printf("[UserService] Failed to delete user ID %s: %v", id.String(), err)
		return err
	}

	log.Printf("[UserService] User deleted successfully: ID %s", id.String())
	return nil
}

func (s *userService) GetEligibleUsers(ctx context.Context) ([]*models.User, error) {
	log.Printf("[UserService] Getting all eligible users")

	users, err := s.userRepo.GetEligibleUsers(ctx)
	if err != nil {
		log.Printf("[UserService] Failed to get eligible users: %v", err)
		return nil, err
	}

	log.Printf("[UserService] Found %d eligible users", len(users))
	return users, nil
}

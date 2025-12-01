package repositories

import (
	"context"
	"fmt"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, user *models.CreateUserRequest) (*models.User, error)
	Update(ctx context.Context, id uuid.UUID, user *models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	IsEligible(ctx context.Context, phone string) (bool, error)
	GetEligibleUsers(ctx context.Context) ([]*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT id, name, phone, email, is_active, created_at, updated_at 
		FROM users 
		WHERE phone = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.Name, &user.Phone, &user.Email,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, phone, email, is_active, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Phone, &user.Email,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	query := `
		INSERT INTO users (name, phone, email) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, phone, email, is_active, created_at, updated_at
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, req.Name, req.Phone, req.Email).Scan(
		&user.ID, &user.Name, &user.Phone, &user.Email,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.Name)
		argIndex++
	}

	if req.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, req.Email)
		argIndex++
	}

	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id) // No updates, return existing user
	}

	// Always update updated_at
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")

	query := fmt.Sprintf(`
		UPDATE users 
		SET %s
		WHERE id = $%d 
		RETURNING id, name, phone, email, is_active, created_at, updated_at
	`, fmt.Sprintf("%s", setParts), argIndex)

	args = append(args, id)

	var user models.User
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Name, &user.Phone, &user.Email,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) IsEligible(ctx context.Context, phone string) (bool, error) {
	query := `SELECT is_active FROM users WHERE phone = $1`

	var isActive bool
	err := r.db.QueryRow(ctx, query, phone).Scan(&isActive)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil // User not found means not eligible
		}
		return false, fmt.Errorf("failed to check user eligibility: %w", err)
	}

	return isActive, nil
}

func (r *userRepository) GetEligibleUsers(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, name, phone, email, is_active, created_at, updated_at 
		FROM users 
		WHERE is_active = true
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get eligible users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Phone, &user.Email,
			&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}

	return users, nil
}

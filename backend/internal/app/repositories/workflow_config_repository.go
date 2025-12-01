package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkflowConfigRepository interface {
	GetActiveWorkflowType(ctx context.Context) (string, error)
}

type workflowConfigRepository struct {
	db *pgxpool.Pool
}

func NewWorkflowConfigRepository(db *pgxpool.Pool) WorkflowConfigRepository {
	return &workflowConfigRepository{db: db}
}

func (r *workflowConfigRepository) GetActiveWorkflowType(ctx context.Context) (string, error) {
	query := `
		SELECT workflow_type 
		FROM workflow_config 
		WHERE id = 1 AND is_active = true
	`

	var workflowType string
	err := r.db.QueryRow(ctx, query).Scan(&workflowType)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("no active workflow configuration found")
		}
		return "", fmt.Errorf("failed to get active workflow type: %w", err)
	}

	return workflowType, nil
}

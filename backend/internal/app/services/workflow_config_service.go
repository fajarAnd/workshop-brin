package services

import (
	"context"
	"log"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/repositories"
)

type WorkflowConfigService interface {
	GetActiveWorkflowType(ctx context.Context) (string, error)
}

type workflowConfigService struct {
	workflowConfigRepo repositories.WorkflowConfigRepository
}

func NewWorkflowConfigService(workflowConfigRepo repositories.WorkflowConfigRepository) WorkflowConfigService {
	return &workflowConfigService{
		workflowConfigRepo: workflowConfigRepo,
	}
}

func (s *workflowConfigService) GetActiveWorkflowType(ctx context.Context) (string, error) {
	log.Printf("[WorkflowConfigService] Getting active workflow type")

	workflowType, err := s.workflowConfigRepo.GetActiveWorkflowType(ctx)
	if err != nil {
		log.Printf("[WorkflowConfigService] Failed to get active workflow type: %v", err)
		return "", err
	}

	log.Printf("[WorkflowConfigService] Active workflow type: %s", workflowType)
	return workflowType, nil
}

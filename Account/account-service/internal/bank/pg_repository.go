package bank

import (
	"2margin.vn/account-service/internal/dto"
	"2margin.vn/account-service/internal/models"
	"context"
)

type PGRepository interface {
	Create(ctx context.Context, bank *models.Bank) (*models.Bank, error)
	Search(ctx context.Context, search *dto.SearchBank) ([]*models.Bank, error)
	FindById(ctx context.Context, bankId string) (*models.Bank, error)
	DeleteById(ctx context.Context, bankId string) (*string, error)
}

package repository

import (
	"acbs.com.vn/account-service/internal/models"
	"context"
)

type AccountRepositoryMock struct {
	CreateFunc func(ctx context.Context, account *models.Account) (*models.Account, error)
}

func (m *AccountRepositoryMock) Create(ctx context.Context, account *models.Account) (*models.Account, error) {
	return m.CreateFunc(ctx, account)
}

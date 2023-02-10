package bank

import (
	"2margin.vn/account-service/internal/dto"
	"context"
)

type UseCase interface {
	CreateBank(ctx context.Context, bank *dto.CreateBank) (string, error)
	GetBankByID(ctx context.Context, bankId string) (*dto.Bank, error)
	DeleteBankByID(ctx context.Context, bankId string) error
}

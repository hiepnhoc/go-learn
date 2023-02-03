package queries

import (
	"acbs.com.vn/account-service/config"
	"acbs.com.vn/account-service/internal/account/repository"
	"acbs.com.vn/account-service/internal/models"
	"acbs.com.vn/account-service/pkg/logger"
	"context"
)

type GetAccountByNameHandler interface {
	Handle(ctx context.Context, query *GetAccountByNameQuery) ([]*models.Account, error)
}

type getAccountByNameHandler struct {
	log         logger.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
}

func NewGetAccountByIdHandler(log logger.Logger, cfg *config.Config, accountRepo repository.AccountRepository) *getAccountByNameHandler {
	return &getAccountByNameHandler{log: log, cfg: cfg, accountRepo: accountRepo}
}

func (q *getAccountByNameHandler) Handle(ctx context.Context, query *GetAccountByNameQuery) ([]*models.Account, error) {
	return q.accountRepo.List(ctx, query.Name)
}

package commands

import (
	"acbs.com.vn/account-service/config"
	"acbs.com.vn/account-service/internal/account/repository"
	"acbs.com.vn/account-service/internal/models"
	"acbs.com.vn/account-service/pkg/logger"
	"context"
	"github.com/opentracing/opentracing-go"
)

type CreateAccountCmdHandler interface {
	Handle(ctx context.Context, command *CreateAccountCommand) (*models.Account, error)
}

type createAccountHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo repository.AccountRepository
}

func NewCreateAccountHandler(log logger.Logger, cfg *config.Config, pgRepo repository.AccountRepository) *createAccountHandler {
	return &createAccountHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *createAccountHandler) Handle(ctx context.Context, command *CreateAccountCommand) (*models.Account, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createAccountHandler.Handle")
	defer span.Finish()

	productDto := &models.Account{Name: command.Name}

	product, err := c.pgRepo.Create(ctx, productDto)
	if err != nil {
		return nil, err
	}

	return product, err
}

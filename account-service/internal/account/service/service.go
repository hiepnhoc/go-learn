package service

import (
	"acbs.com.vn/account-service/config"
	"acbs.com.vn/account-service/internal/account/commands"
	"acbs.com.vn/account-service/internal/account/queries"
	"acbs.com.vn/account-service/internal/account/repository"
	"acbs.com.vn/account-service/pkg/logger"
)

type AccountService struct {
	Commands *commands.AccountCommands
	Queries  *queries.AccountQueries
}

func NewAccountService(
	log logger.Logger,
	cfg *config.Config,
	pgRepo repository.AccountRepository,
) *AccountService {

	createAccountHandler := commands.NewCreateAccountHandler(log, cfg, pgRepo)
	getAccountByNameHandler := queries.NewGetAccountByIdHandler(log, cfg, pgRepo)

	createAccountCommands := commands.NewAccountCommands(createAccountHandler)
	accountQueries := queries.NewAccountQueries(getAccountByNameHandler)

	return &AccountService{Commands: createAccountCommands, Queries: accountQueries}
}

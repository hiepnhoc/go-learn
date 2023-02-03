package service

import (
	"acbs.com.vn/account-service/config"
	commands2 "acbs.com.vn/account-service/internal/account/commands"
	queries2 "acbs.com.vn/account-service/internal/account/queries"
	"acbs.com.vn/account-service/internal/account/repository"
	"acbs.com.vn/account-service/pkg/logger"
)

type AccountService struct {
	Commands *commands2.AccountCommands
	Queries  *queries2.AccountQueries
}

func NewAccountService(
	log logger.Logger,
	cfg *config.Config,
	pgRepo repository.AccountRepository,
) *AccountService {

	createAccountHandler := commands2.NewCreateAccountHandler(log, cfg, pgRepo)
	getAccountByNameHandler := queries2.NewGetAccountByIdHandler(log, cfg, pgRepo)

	createAccountCommands := commands2.NewAccountCommands(createAccountHandler)
	accountQueries := queries2.NewAccountQueries(getAccountByNameHandler)

	return &AccountService{Commands: createAccountCommands, Queries: accountQueries}
}

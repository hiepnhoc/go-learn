package repository

import (
	"acbs.com.vn/account-service/config"
	"acbs.com.vn/account-service/internal/models"
	"acbs.com.vn/account-service/pkg/logger"
	"context"
	uuid "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	idPrefix = "ac-i-"
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) (*models.Account, error)
	List(ctx context.Context, name string) ([]*models.Account, error)
}

type accountRepository struct {
	log logger.Logger
	cfg *config.Config
	db  *sqlx.DB
}

func NewAccountRepository(log logger.Logger, cfg *config.Config, db *sqlx.DB) *accountRepository {
	return &accountRepository{log: log, cfg: cfg, db: db}
}

const (
	InsertIntoAccounts = "INSERT INTO accounts (id, name) VALUES ($1,$2) RETURNING id, name"
	SearchAccounts     = "SELECT * FROM accounts where name like '%' || $1 || '%'"
)

func (r *accountRepository) Create(ctx context.Context, req *models.Account) (*models.Account, error) {

	id := idPrefix + uuid.New().String()

	var created models.Account
	if err := r.db.QueryRowContext(ctx, InsertIntoAccounts, id, req.Name).Scan(&created.Id, &created.Name); err != nil {
		return nil, errors.New(err.Error())
	}

	return &created, nil
}

func (r *accountRepository) List(ctx context.Context, name string) ([]*models.Account, error) {
	if name == "" {
		return nil, errors.New("filter cannot be nil")
	}

	var accounts []*models.Account
	err := r.db.SelectContext(ctx, &accounts, SearchAccounts, name)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return accounts, nil

}

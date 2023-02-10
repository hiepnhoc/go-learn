package repository

import (
	"2margin.vn/account-service/internal/dto"
	"2margin.vn/account-service/internal/models"
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	TableName      = "bank"
	IdPrefix       = TableName + "-"
	BankIdColumn   = "id"
	BankNameColumn = "bank_name"
	BankCodeColumn = "bank_code"
)

type pgRepository struct {
	table string
	db    *sqlx.DB
}

func NewPGRepository(db *sqlx.DB) *pgRepository {
	return &pgRepository{table: TableName, db: db}
}

func (r *pgRepository) Create(ctx context.Context, bank *models.Bank) (*models.Bank, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankPgRepository.Create")
	defer span.Finish()

	var bankId string

	if len(strings.TrimSpace(bank.Id)) == 0 {
		bank.Id = IdPrefix + uuid.New().String()
	}

	sql, _, _ := goqu.Insert(r.table).Rows(bank).Returning(BankIdColumn).ToSQL()

	if err := r.db.QueryRowContext(ctx, sql).Scan(&bankId); err != nil {
		return nil, errors.New("BankPgRepository.Create.QueryRowContext " + err.Error())
	}

	if bankId != bank.Id {
		return nil, errors.Errorf("BankPgRepository.Create CheckBankId Want %s, Got %s "+bank.Id, bankId)
	}

	return r.FindById(ctx, bankId)

}

func (r *pgRepository) Search(ctx context.Context, search *dto.SearchBank) ([]*models.Bank, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "BankPgRepository.FindById")
	defer span.Finish()

	var banks []*models.Bank

	var where = goqu.And()

	if search.BankId != nil {
		where = where.Append(goqu.C(BankIdColumn).Eq(search.BankId))
	}

	if search.BankName != nil {
		where = where.Append(goqu.C(BankNameColumn).Eq(search.BankName))
	}

	if search.BankCode != nil {
		where = where.Append(goqu.C(BankCodeColumn).Eq(search.BankCode))
	}

	sql, _, _ := goqu.From(r.table).Select().Where(where).ToSQL()

	if err := r.db.SelectContext(ctx, &banks, sql); err != nil {
		return nil, errors.Wrap(err, "BankPgRepository.Search.SelectContext")
	}

	return banks, nil
}

func (r *pgRepository) FindById(ctx context.Context, bankId string) (*models.Bank, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "BankPgRepository.FindById")
	defer span.Finish()

	bank := &models.Bank{}

	var where = goqu.And(
		goqu.C(BankIdColumn).Eq(bankId),
	)

	sql, _, _ := goqu.From(r.table).Select().Where(where).Limit(1).ToSQL()

	if err := r.db.GetContext(ctx, bank, sql); err != nil {
		return nil, errors.Wrap(err, "BankPgRepository.FindById.GetContext")
	}

	return bank, nil
}

func (r *pgRepository) DeleteById(ctx context.Context, bankId string) (*string, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "BankPgRepository.DeleteById")

	defer span.Finish()

	var deleteId string

	var where = goqu.And(
		goqu.C(BankIdColumn).Eq(bankId),
	)

	type update struct {
		DeletedAt time.Time `db:"deleted_at"`
	}

	sql, _, _ := goqu.Update(r.table).Where(where).Set(&update{DeletedAt: time.Now()}).Returning(BankIdColumn).ToSQL()

	if err := r.db.QueryRowContext(ctx, sql).Scan(&deleteId); err != nil {
		return nil, errors.New("BankPgRepository.DeleteById.QueryRowContext " + err.Error())
	}

	if bankId != deleteId {
		return nil, errors.Errorf("BankPgRepository.DeleteById CheckBankId Want %s, Got %s "+bankId, deleteId)
	}

	return &deleteId, nil
}

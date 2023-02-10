package repository

import (
	"2margin.vn/account-service/internal/dto"
	"2margin.vn/account-service/internal/models"
	"2margin.vn/account-service/pkg/postgres"
	testPostgres "2margin.vn/account-service/test/containers/postgres"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

var container *testPostgres.Container

func TestMain(m *testing.M) {

	container = testPostgres.NewTestContainer(m)
	err := container.Run()
	if err != nil {
		return
	}

}

func TestPgRepository_Create_FindById_Delete(t *testing.T) {

	type args struct {
		ctx context.Context
		req *models.Bank
	}
	type want struct {
		res *models.Bank
		err error
	}

	type test struct {
		name string
		args
		want
	}

	bank1 := newRandomBank()
	bank2 := newRandomBank()
	bank3 := newRandomBank()
	bank3.Id = uuid.NewString()

	tests := []test{
		{
			name: bank1.BankName,
			args: args{
				ctx: context.Background(),
				req: bank1,
			},
			want: want{
				res: bank1,
				err: nil,
			},
		},
		{
			name: bank2.BankName,
			args: args{
				ctx: context.Background(),
				req: bank2,
			},
			want: want{
				res: bank2,
				err: nil,
			},
		},
		{
			name: bank3.BankName,
			args: args{
				ctx: context.Background(),
				req: bank3,
			},
			want: want{
				res: bank3,
				err: nil,
			},
		},
	}

	sqlxDB, err := postgres.NewSqlxDB(container.Cfg)
	if err != nil {
		t.Fatalf("Initial sqlx failed %v", err)
	}

	pgRepository := NewPGRepository(sqlxDB)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := pgRepository.Create(test.args.ctx, test.args.req)
			assert.Equal(t, test.want.err, err)
			if err == nil {
				assert.Equal(t, test.want.res.Id, got.Id)
				assert.Equal(t, test.want.res.BankCode, got.BankCode)
				assert.Equal(t, test.want.res.BankName, got.BankName)
				assert.Equal(t, test.want.res.BankNumber, got.BankNumber)
				assert.Equal(t, test.want.res.CustomerName, got.CustomerName)
				assert.Equal(t, test.want.res.Content, got.Content)
				assert.Equal(t, test.want.res.IsDefault, got.IsDefault)
				assert.Equal(t, test.want.res.Owner, got.Owner)
			}
		})
	}

	t.Run("DELETE "+bank3.BankName, func(t *testing.T) {
		var ctx = context.Background()
		deleteId, err := pgRepository.DeleteById(ctx, bank3.Id)
		if err == nil {
			assert.Equal(t, bank3.Id, *deleteId)
		}
		bank, err := pgRepository.FindById(ctx, bank3.Id)
		if err != nil {
			t.Error(err)
		}
		assert.NotEmpty(t, bank.DeletedAt)

	})

}

func TestPgRepository_Search(t *testing.T) {

	type args struct {
		ctx context.Context
		req *dto.SearchBank
	}
	type want struct {
		res []*models.Bank
		err error
	}

	type test struct {
		name string
		args
		want
	}

	bank1 := newRandomBank()

	bank2 := newRandomBank()

	bank3 := newRandomBank()

	var code2 = "BANK_SEARCH_2"
	bank2.BankCode = code2
	bank3.BankCode = code2

	search1 := &dto.SearchBank{
		BankId: &bank1.Id,
	}

	search2 := &dto.SearchBank{
		BankCode: &code2,
	}

	search3 := &dto.SearchBank{
		BankId:   &bank1.Id,
		BankCode: &code2,
	}

	res1 := []*models.Bank{bank1}

	res2 := []*models.Bank{bank2, bank3}

	tests := []test{
		{
			name: "SEARCH_BANK_ID_1",
			args: args{
				ctx: context.Background(),
				req: search1,
			},
			want: want{
				res: res1,
				err: nil,
			},
		},
		{
			name: "SEARCH_BANK_CODE",
			args: args{
				ctx: context.Background(),
				req: search2,
			},
			want: want{
				res: res2,
				err: nil,
			},
		},

		{
			name: "SEARCH_3",
			args: args{
				ctx: context.Background(),
				req: search3,
			},
			want: want{
				res: []*models.Bank{},
				err: nil,
			},
		},
	}

	sqlxDB, err := postgres.NewSqlxDB(container.Cfg)
	if err != nil {
		t.Fatalf("Initial sqlx failed %v", err)
	}

	pgRepository := NewPGRepository(sqlxDB)

	_, _ = pgRepository.Create(context.Background(), bank1)
	_, _ = pgRepository.Create(context.Background(), bank2)
	_, _ = pgRepository.Create(context.Background(), bank3)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := pgRepository.Search(test.args.ctx, test.args.req)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, len(test.want.res), len(got))

			for _, bankGot := range got {
				var pass = false
				for _, bankRes := range test.want.res {
					if bankGot.Id == bankRes.Id {
						pass = true
					}
				}
				if pass == false {
					t.Fatalf("Search result isnt correct")
				}
			}
		})
	}

}

func newRandomBank() *models.Bank {
	bank := &models.Bank{}

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)

	var content = "BANK_CONTENT_" + suffix
	bank.BankCode = "BANK_" + suffix
	bank.BankName = "BANK_NAME_" + suffix
	bank.BankNumber = "BANK_NUMBER_" + suffix
	bank.CustomerName = "BANK_CUSTOMER_" + suffix
	bank.Content = &content
	bank.IsDefault = false
	bank.Owner = "BANK_OWNER_" + suffix

	return bank
}

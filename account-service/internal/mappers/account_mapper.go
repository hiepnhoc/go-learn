package mappers

import (
	"acbs.com.vn/account-service/internal/models"
	"acbs.com.vn/account-service/pkg/utils"
	accountService "acbs.com.vn/account-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AccountsList list response with pagination
type AccountsList struct {
	TotalCount int64             `json:"totalCount" bson:"totalCount"`
	TotalPages int64             `json:"totalPages" bson:"totalPages"`
	Page       int64             `json:"page" bson:"page"`
	Size       int64             `json:"size" bson:"size"`
	HasMore    bool              `json:"hasMore" bson:"hasMore"`
	Accounts   []*models.Account `json:"products" bson:"products"`
}

func NewAccountListWithPagination(accounts []*models.Account, count int64, pagination *utils.Pagination) *AccountsList {
	return &AccountsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Accounts:   accounts,
	}
}

func AccountToGrpcMessage(account *models.Account) *accountService.Account {
	return &accountService.Account{
		Id:        account.Id,
		Name:      account.Name,
		CreatedAt: timestamppb.New(account.CreatedAt),
		UpdatedAt: timestamppb.New(account.UpdatedAt),
	}
}

func AccountListToGrpc(accounts *AccountsList) *accountService.SearchRes {
	list := make([]*accountService.Account, 0, len(accounts.Accounts))
	for _, account := range accounts.Accounts {
		list = append(list, AccountToGrpcMessage(account))
	}

	return &accountService.SearchRes{
		TotalCount: accounts.TotalCount,
		TotalPages: accounts.TotalPages,
		Page:       accounts.Page,
		Size:       accounts.Size,
		HasMore:    accounts.HasMore,
		Accounts:   list,
	}
}

func AccountListToAccounts(accounts []*models.Account) []*accountService.Account {
	list := make([]*accountService.Account, 0, len(accounts))
	for _, account := range accounts {
		list = append(list, AccountToGrpcMessage(account))
	}

	return list
}

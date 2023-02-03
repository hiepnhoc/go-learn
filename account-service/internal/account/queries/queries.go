package queries

type AccountQueries struct {
	GetAccountByName GetAccountByNameHandler
}

func NewAccountQueries(getAccountByName GetAccountByNameHandler) *AccountQueries {
	return &AccountQueries{GetAccountByName: getAccountByName}
}

type GetAccountByNameQuery struct {
	Name string `json:"name" validate:"required,gte=0,lte=255"`
}

func NewGetAccountByNameQuery(name string) *GetAccountByNameQuery {
	return &GetAccountByNameQuery{Name: name}
}

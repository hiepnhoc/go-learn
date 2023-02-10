package dto

type CreateBank struct {
	BankCode     string `validate:"required"`
	BankName     string `validate:"required"`
	BankBranch   string `validate:"required"`
	CustomerName string `validate:"required"`
}

type SearchBank struct {
	BankId       *string
	BankCode     *string
	BankName     *string
	CustomerName *string
	SkipDeleted  *bool
}

type Bank struct {
	BankCode     string `validate:"required"`
	BankName     string `validate:"required"`
	CustomerName string `validate:"required"`
}

package commands

type AccountCommands struct {
	CreateAccount CreateAccountCmdHandler
}

func NewAccountCommands(createAccount CreateAccountCmdHandler) *AccountCommands {
	return &AccountCommands{CreateAccount: createAccount}
}

type CreateAccountCommand struct {
	Name string `json:"name" validate:"required,gte=0,lte=255"`
}

func NewCreateAccountCommand(name string) *CreateAccountCommand {
	return &CreateAccountCommand{Name: name}
}

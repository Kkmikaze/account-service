package schema

type CreateAccountRequest struct {
	CitizenID string `json:"citizen_id" validate:"required,len=16"`
	Name      string `json:"name" validate:"required,min=3,max=150"`
	Phone     string `json:"phone" validate:"required,e164"`
}

type SavingBalanceRequest struct {
	AccountNumber string  `json:"account_number" validate:"required,len=14"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

type WithdrawBalanceRequest struct {
	AccountNumber string  `json:"account_number" validate:"required,len=14"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

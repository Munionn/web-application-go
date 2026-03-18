package types

type CreateAccountRequest struct {
	Name string `json:"name"`
}

type UpdateAccountRequest struct {
	Name         *string  `json:"name"`
	BaseCurrence *string  `json:"base_currence"`
	Balance      *float64 `json:"balance"`
}

// CreateTransactionRequest is the body for creating a transaction (e.g. POST /accounts/:id/transactions)
type CreateTransactionRequest struct {
	Amount           uint   `json:"amount"`
	BaseCurrency     string `json:"base_currency"`
	Type             string `json:"type"`
	ShortDescription string `json:"short_description"`
}

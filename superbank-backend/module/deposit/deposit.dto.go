package deposit

type CreateDepositRequest struct {
	CustomerID   string  `json:"customer_id" validate:"required,uuid4"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	InterestRate float64 `json:"interest_rate" validate:"required,gte=0"`
	TermMonths   int     `json:"term_months" validate:"required,gt=0"`
	StartDate    string  `json:"start_date" validate:"required,datetime=2006-01-02"`
	Note         string  `json:"note"`
}

package pocket

type CreatePocketRequest struct {
	CustomerID   string   `json:"customer_id" validate:"required,uuid4"`
	Name         string   `json:"name" validate:"required,min=3"`
	TargetAmount *float64 `json:"targetAmount,omitempty"`
	TargetDate   *string  `json:"targetDate,omitempty"`
}

type TopUpOrWithdrawPocketRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

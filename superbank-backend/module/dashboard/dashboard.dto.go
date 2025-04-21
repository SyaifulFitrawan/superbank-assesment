package dashboard

type DashboardTotalCounts struct {
	TotalCustomers int64 `json:"total_customers"`
	TotalDeposits  int64 `json:"total_deposits"`
	TotalPockets   int64 `json:"total_pockets"`
}

type AccountType struct {
	AccountType string `json:"account_type"`
	Count       int    `json:"count"`
}

type CustomerDepositOrPocketGroup struct {
	RangeLabel string `json:"range_label"`
	Count      int64  `json:"count"`
}

type DashboardResponse struct {
	Total   DashboardTotalCounts           `json:"total"`
	Type    []AccountType                  `json:"type"`
	Deposit []CustomerDepositOrPocketGroup `json:"deposits"`
	Pocket  []CustomerDepositOrPocketGroup `json:"pockets"`
}

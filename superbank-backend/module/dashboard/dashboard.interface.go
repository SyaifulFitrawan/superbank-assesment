package dashboard

import "context"

type DashboardRepository interface {
	GetTotals(ctx context.Context) (DashboardTotalCounts, error)
	CountByAccountType(ctx context.Context) ([]AccountType, error)
	GetCustomerDepositGroups(ctx context.Context) ([]CustomerDepositOrPocketGroup, error)
	GetCustomerPocketGroups(ctx context.Context) ([]CustomerDepositOrPocketGroup, error)
}

type DashboardService interface {
	GetDashboard(ctx context.Context) (DashboardResponse, error)
}

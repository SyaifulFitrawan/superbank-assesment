package dashboard

import (
	"context"
)

type dashboardServiceImpl struct {
	repo DashboardRepository
}

func NewDashboardService(repo DashboardRepository) DashboardService {
	return &dashboardServiceImpl{repo: repo}
}

func (s *dashboardServiceImpl) GetDashboard(ctx context.Context) (DashboardResponse, error) {
	total, err := s.repo.GetTotals(ctx)
	if err != nil {
		return DashboardResponse{}, err
	}

	countTyoe, err := s.repo.CountByAccountType(ctx)
	if err != nil {
		return DashboardResponse{}, err
	}

	deposits, err := s.repo.GetCustomerDepositGroups(ctx)
	if err != nil {
		return DashboardResponse{}, err
	}

	pockets, err := s.repo.GetCustomerPocketGroups(ctx)
	if err != nil {
		return DashboardResponse{}, err
	}

	result := DashboardResponse{
		Total:   total,
		Type:    countTyoe,
		Deposit: deposits,
		Pocket:  pockets,
	}

	return result, nil
}

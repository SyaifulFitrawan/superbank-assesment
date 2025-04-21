package deposit

import (
	"bank-backend/model"
	"context"
)

type DepositRepository interface {
	Create(ctx context.Context, input *model.Deposit) error
	Update(ctx context.Context, id string, input *model.Deposit) error
	FindMatureUnwithdraw(ctx context.Context) ([]model.Deposit, error)
}

type DepositService interface {
	Create(ctx context.Context, input CreateDepositRequest) (*model.Deposit, error)
	ProcessMatureDeposits(ctx context.Context) error
}

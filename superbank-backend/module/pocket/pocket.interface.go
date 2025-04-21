package pocket

import (
	"bank-backend/model"
	"context"
)

type PocketRepository interface {
	Create(ctx context.Context, input *model.Pocket) error
	Detail(ctx context.Context, id string) (*model.Pocket, error)
	Update(ctx context.Context, id string, input *model.Pocket) error
	Deactivated(ctx context.Context, id string, input any) error
}

type PocketService interface {
	Create(ctx context.Context, input CreatePocketRequest) (*model.Pocket, error)
	TopUp(ctx context.Context, id string, input TopUpOrWithdrawPocketRequest) error
	Withdrawn(ctx context.Context, id string, input TopUpOrWithdrawPocketRequest) error
	Deactivated(ctx context.Context, id string) error
}

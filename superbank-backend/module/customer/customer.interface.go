package customer

import (
	"bank-backend/model"
	"bank-backend/utils"
	"context"
)

type CustomerRepository interface {
	Create(ctx context.Context, input *model.Customer) error
	List(ctx context.Context, limit, offset int, search string) ([]model.Customer, int64, error)
	Detail(ctx context.Context, id string) (*model.Customer, error)
	Update(ctx context.Context, id string, input *model.Customer) error
	Delete(ctx context.Context, id string) error
	AddBalance(ctx context.Context, customerID string, amount float64) error
}

type CustomerService interface {
	Create(ctx context.Context, input CustomerCreateRequest) (*model.Customer, error)
	List(ctx context.Context, page, limit int, search string) ([]model.Customer, *utils.Paginator, error)
	Detail(ctx context.Context, id string) (*CustomerDetailResponse, error)
	Update(ctx context.Context, id string, input CustomerUpdateRequest) error
	Delete(ctx context.Context, id string) error
}

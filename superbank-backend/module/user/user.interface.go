package user

import (
	"bank-backend/model"
	"bank-backend/utils"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	List(ctx context.Context, limit, offset int, search string) ([]model.User, int64, error)
	Detail(ctx context.Context, id string) (*model.User, error)
	Delete(ctx context.Context, id string) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserService interface {
	Create(ctx context.Context, input *model.User) (*model.User, error)
	List(ctx context.Context, page, limit int, search string) ([]model.User, *utils.Paginator, error)
	Detail(ctx context.Context, id string) (*model.User, error)
	Delete(ctx context.Context, id string) error
}

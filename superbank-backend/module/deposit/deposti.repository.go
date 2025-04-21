package deposit

import (
	"bank-backend/database"
	"bank-backend/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type depositRepositoryImpl struct {
	db *gorm.DB
}

func NewDepositRepository(db *gorm.DB) DepositRepository {
	return &depositRepositoryImpl{db: db}
}

func (r *depositRepositoryImpl) Create(ctx context.Context, input *model.Deposit) error {
	trx := database.FromContext(ctx, r.db)
	return trx.Create(input).Error
}

func (r *depositRepositoryImpl) Update(ctx context.Context, id string, input *model.Deposit) error {
	trx := database.FromContext(ctx, r.db)
	return trx.WithContext(ctx).Model(&model.Deposit{}).Where("id = ?", id).Updates(input).Error
}

func (r *depositRepositoryImpl) FindMatureUnwithdraw(ctx context.Context) ([]model.Deposit, error) {
	var deposits []model.Deposit

	err := r.db.WithContext(ctx).
		Where("is_withdrawn = ? AND maturity_date <= ?", false, time.Now()).
		Find(&deposits).Error

	if err != nil {
		return nil, err
	}

	return deposits, nil
}

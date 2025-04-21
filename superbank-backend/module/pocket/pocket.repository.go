package pocket

import (
	"bank-backend/database"
	"bank-backend/model"
	"context"

	"gorm.io/gorm"
)

type pocketRepositoryImpl struct {
	db *gorm.DB
}

func NewPocketRepository(db *gorm.DB) PocketRepository {
	return &pocketRepositoryImpl{db: db}
}

func (r *pocketRepositoryImpl) Create(ctx context.Context, input *model.Pocket) error {
	return r.db.WithContext(ctx).Create(input).Error
}

func (r *pocketRepositoryImpl) Detail(ctx context.Context, id string) (*model.Pocket, error) {
	var pocket model.Pocket
	if err := r.db.WithContext(ctx).First(&pocket, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &pocket, nil
}

func (r *pocketRepositoryImpl) Update(ctx context.Context, id string, input *model.Pocket) error {
	trx := database.FromContext(ctx, r.db)
	return trx.WithContext(ctx).Model(&model.Pocket{}).Where("id = ?", id).Updates(input).Error
}

func (r *pocketRepositoryImpl) Deactivated(ctx context.Context, id string, input any) error {
	trx := database.FromContext(ctx, r.db)
	return trx.WithContext(ctx).Model(&model.Pocket{}).Where("id = ?", id).Updates(input).Error
}

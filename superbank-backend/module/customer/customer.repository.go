package customer

import (
	"bank-backend/database"
	"bank-backend/model"
	"context"

	"gorm.io/gorm"
)

type customerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepositoryImpl{db: db}
}

func (r *customerRepositoryImpl) Create(ctx context.Context, input *model.Customer) error {
	return r.db.WithContext(ctx).Create(input).Error
}

func (r *customerRepositoryImpl) List(ctx context.Context, limit, offset int, search string) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Customer{})

	if search != "" {
		query = query.Where("account_number ILIKE ? OR name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *customerRepositoryImpl) Detail(ctx context.Context, id string) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.WithContext(ctx).
		Preload("Deposits").
		Preload("Pockets").
		First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepositoryImpl) Update(ctx context.Context, id string, input *model.Customer) error {
	trx := database.FromContext(ctx, r.db)
	return trx.WithContext(ctx).Model(&model.Customer{}).Where("id = ?", id).Updates(input).Error
}

func (r *customerRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Customer{}).Error
}

func (r *customerRepositoryImpl) AddBalance(ctx context.Context, customerID string, amount float64) error {
	trx := database.FromContext(ctx, r.db)
	return trx.Model(&model.Customer{}).
		Where("id = ?", customerID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

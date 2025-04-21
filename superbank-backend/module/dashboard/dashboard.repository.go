package dashboard

import (
	"bank-backend/model"
	"context"

	"gorm.io/gorm"
)

type dashboardRepositoryImpl struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepositoryImpl{db: db}
}

func (r *dashboardRepositoryImpl) GetTotals(ctx context.Context) (DashboardTotalCounts, error) {
	var result DashboardTotalCounts

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				(SELECT COUNT(*) FROM customers) AS total_customers,
				(SELECT COUNT(*) FROM deposits) AS total_deposits,
				(SELECT COUNT(*) FROM pockets) AS total_pockets
		`).Scan(&result).Error

	return result, err
}

func (r *dashboardRepositoryImpl) CountByAccountType(ctx context.Context) ([]AccountType, error) {
	var results []AccountType

	err := r.db.WithContext(ctx).
		Model(&model.Customer{}).
		Select("account_type, COUNT(*) as count").
		Group("account_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *dashboardRepositoryImpl) GetCustomerDepositGroups(ctx context.Context) ([]CustomerDepositOrPocketGroup, error) {
	var result []CustomerDepositOrPocketGroup

	err := r.db.WithContext(ctx).Raw(`
		WITH ranges AS (
			SELECT '0-1' AS range_label
			UNION SELECT '2-3'
			UNION SELECT '4-5'
			UNION SELECT '6+'
		),
		deposit_counts AS (
			SELECT
				CASE
					WHEN d.total BETWEEN 0 AND 1 THEN '0-1'
					WHEN d.total BETWEEN 2 AND 3 THEN '2-3'
					WHEN d.total BETWEEN 4 AND 5 THEN '4-5'
					ELSE '6+'
				END AS range_label,
				COUNT(*) AS count
			FROM (
				SELECT customer_id, COUNT(*) as total
				FROM deposits
				GROUP BY customer_id
			) d
			GROUP BY range_label
		)
		SELECT r.range_label, COALESCE(dc.count, 0) AS count
		FROM ranges r
		LEFT JOIN deposit_counts dc ON r.range_label = dc.range_label
		ORDER BY r.range_label
	`).Scan(&result).Error

	return result, err
}

func (r *dashboardRepositoryImpl) GetCustomerPocketGroups(ctx context.Context) ([]CustomerDepositOrPocketGroup, error) {
	var result []CustomerDepositOrPocketGroup

	err := r.db.WithContext(ctx).Raw(`
		WITH ranges AS (
			SELECT '0-1' AS range_label
			UNION SELECT '2-3'
			UNION SELECT '4-5'
			UNION SELECT '6+'
		),
		deposit_counts AS (
			SELECT
				CASE
					WHEN p.total BETWEEN 0 AND 1 THEN '0-1'
					WHEN p.total BETWEEN 2 AND 3 THEN '2-3'
					WHEN p.total BETWEEN 4 AND 5 THEN '4-5'
					ELSE '6+'
				END AS range_label,
				COUNT(*) AS count
			FROM (
				SELECT customer_id, COUNT(*) as total
				FROM pockets
				GROUP BY customer_id
			) p
			GROUP BY range_label
		)
		SELECT r.range_label, COALESCE(dc.count, 0) AS count
		FROM ranges r
		LEFT JOIN deposit_counts dc ON r.range_label = dc.range_label
		ORDER BY r.range_label
	`).Scan(&result).Error

	return result, err
}

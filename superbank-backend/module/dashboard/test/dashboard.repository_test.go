package dashboard

import (
	"bank-backend/module/dashboard"
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestGetTotals(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := dashboard.NewDashboardRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			(SELECT COUNT(*) FROM customers) AS total_customers,
			(SELECT COUNT(*) FROM deposits) AS total_deposits,
			(SELECT COUNT(*) FROM pockets) AS total_pockets
	`)).
		WillReturnRows(sqlmock.NewRows([]string{"total_customers", "total_deposits", "total_pockets"}).
			AddRow(int64(5), int64(10), int64(7)))

	result, err := repo.GetTotals(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(5), result.TotalCustomers)
	require.Equal(t, int64(10), result.TotalDeposits)
	require.Equal(t, int64(7), result.TotalPockets)
}

func TestCountByAccountType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := dashboard.NewDashboardRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT account_type, COUNT(*) as count FROM "customers" GROUP BY "account_type"`,
		)).WillReturnRows(sqlmock.NewRows([]string{"account_type", "count"}).
			AddRow("basic", int(3)).
			AddRow("premium", int64(2)))

		result, err := repo.CountByAccountType(context.Background())
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, "basic", result[0].AccountType)
		require.Equal(t, int(3), result[0].Count)
	})

	t.Run("query error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := dashboard.NewDashboardRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT account_type, COUNT(*) as count FROM "customers" GROUP BY "account_type"`,
		)).WillReturnError(fmt.Errorf("db error"))

		result, err := repo.CountByAccountType(context.Background())
		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestGetCustomerDepositGroups(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := dashboard.NewDashboardRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).WillReturnRows(sqlmock.NewRows([]string{"range_label", "count"}).
		AddRow("0-1", int64(2)).
		AddRow("2-3", int64(3)).
		AddRow("4-5", int64(1)).
		AddRow("6+", int64(0)))

	result, err := repo.GetCustomerDepositGroups(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 4)
	require.Equal(t, "0-1", result[0].RangeLabel)
	require.Equal(t, int64(2), result[0].Count)
}

func TestGetCustomerPocketGroups(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := dashboard.NewDashboardRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).WillReturnRows(sqlmock.NewRows([]string{"range_label", "count"}).
		AddRow("0-1", int64(5)).
		AddRow("2-3", int64(0)).
		AddRow("4-5", int64(1)).
		AddRow("6+", int64(1)))

	result, err := repo.GetCustomerPocketGroups(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 4)
	require.Equal(t, "6+", result[3].RangeLabel)
	require.Equal(t, int64(1), result[3].Count)
}

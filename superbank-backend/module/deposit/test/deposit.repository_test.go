package deposit

import (
	"bank-backend/database"
	"bank-backend/model"
	"bank-backend/module/deposit"
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateDeposit(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := deposit.NewDepositRepository(db)
	ctx := context.Background()

	customer := &model.Deposit{
		CustomerID:   uuid.NewV4(),
		Amount:       100000,
		InterestRate: 5.5,
		TermMonths:   12,
		Note:         "Test",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "deposits" ("id","customer_id","amount","interest_rate","term_months","start_date","maturity_date","is_withdrawn","note","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)).
		WithArgs(
			sqlmock.AnyArg(),
			customer.CustomerID,
			customer.Amount,
			customer.InterestRate,
			customer.TermMonths,
			AnyTime{},
			AnyTime{},
			false,
			customer.Note,
			AnyTime{},
			AnyTime{},
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, customer)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, customer.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateDeposit(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := deposit.NewDepositRepository(gormDB)

	id := uuid.NewV4().String()
	input := &model.Deposit{
		Amount: 100000,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "deposits"`).
		WithArgs(
			input.Amount,
			AnyTime{},
			id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	trx := gormDB.Begin()
	ctx := database.NewContext(context.Background(), trx)

	err = repo.Update(ctx, id, input)
	require.NoError(t, err)
	require.NoError(t, trx.Commit().Error)
}

func TestFindMatureUnwithdraw(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := deposit.NewDepositRepository(db)
	ctx := context.Background()

	t.Run("should return deposit with matutity date is false", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "customer_id", "amount", "interest_rate", "term_months",
			"start_date", "maturity_date", "is_withdrawn", "note", "created_at", "updated_at",
		}).AddRow(
			uuid.NewV4(), uuid.NewV4(), 100000, 5.5, 12,
			now.AddDate(0, -12, 0), now.Add(-time.Hour), false, "Test note", now, now,
		)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "deposits" WHERE is_withdrawn = $1 AND maturity_date <= $2`)).
			WithArgs(false, sqlmock.AnyArg()).
			WillReturnRows(rows)

		result, err := repo.FindMatureUnwithdraw(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.False(t, result[0].IsWithdrawn)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return failed deposit with matutity date is false", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "deposits" WHERE is_withdrawn = $1 AND maturity_date <= $2`)).
			WithArgs(false, sqlmock.AnyArg()).
			WillReturnError(errors.New("query failed"))

		result, err := repo.FindMatureUnwithdraw(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "query failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

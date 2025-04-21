package pocket

import (
	"bank-backend/database"
	"bank-backend/model"
	"bank-backend/module/pocket"
	"context"
	"database/sql/driver"
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

	repo := pocket.NewPocketRepository(db)
	ctx := context.Background()

	customer := &model.Pocket{
		CustomerID: uuid.NewV4(),
		Name:       "test",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "pockets" ("customer_id","name","balance","is_active","created_at","updated_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id","target_amount","target_date"`)).
		WithArgs(
			customer.CustomerID,
			customer.Name,
			float64(0),
			true,
			AnyTime{},
			AnyTime{},
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "target_amount", "target_date"}).
			AddRow(uuid.NewV4(), nil, nil),
		)
	mock.ExpectCommit()

	err := repo.Create(ctx, customer)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, customer.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDetailPocket(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := pocket.NewPocketRepository(db)
	ctx := context.Background()

	pocketId := uuid.NewV4()
	expectedPocket := model.Pocket{
		ID:           pocketId,
		CustomerID:   uuid.NewV4(),
		Name:         "test",
		Balance:      float64(0),
		TargetAmount: nil,
		TargetDate:   nil,
		IsActive:     true,
	}

	t.Run("should return detail of pocket", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "pockets" WHERE id = \$1 ORDER BY "pockets"\."id" LIMIT \$2`).
			WithArgs(pocketId, 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "customer_id", "name", "balance", "target_amount", "target_date", "is_active", "created_at", "updated_at",
			}).AddRow(
				expectedPocket.ID,
				expectedPocket.CustomerID,
				expectedPocket.Name,
				expectedPocket.Balance,
				expectedPocket.TargetAmount,
				expectedPocket.TargetDate,
				expectedPocket.IsActive,
				nil,
				nil,
			))

		result, err := repo.Detail(ctx, pocketId.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedPocket.Name, result.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error if pocket not found", func(t *testing.T) {
		nonExistentID := uuid.NewV4()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pockets" WHERE id = $1 ORDER BY "pockets"."id" LIMIT $2`)).
			WithArgs(nonExistentID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.Detail(ctx, nonExistentID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdatePocket(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := pocket.NewPocketRepository(gormDB)

	id := uuid.NewV4().String()
	input := &model.Pocket{
		Name: "test",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "pockets"`).
		WithArgs(
			input.Name,
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

func TestDeactivatedPocket(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := pocket.NewPocketRepository(gormDB)

	id := uuid.NewV4().String()
	input := &model.Pocket{
		Balance:  0.0,
		IsActive: false,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "pockets" SET "updated_at"=\$1 WHERE id = \$2`).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	trx := gormDB.Begin()
	ctx := database.NewContext(context.Background(), trx)

	err = repo.Deactivated(ctx, id, input)
	require.NoError(t, err)
	require.NoError(t, trx.Commit().Error)
}

package customer

import (
	"bank-backend/database"
	"bank-backend/model"
	"bank-backend/module/customer"
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

func TestCreateCustomer(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := customer.NewCustomerRepository(db)
	ctx := context.Background()

	customer := &model.Customer{
		Name:          "John Doe",
		Phone:         "081234567890",
		Address:       "Jakarta",
		ParentName:    "Jane Doe",
		AccountNumber: "1234567890",
		AccountBranch: "Jakarta",
		AccountType:   "Savings",
		Balance:       1000000,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "customers" ("id","name","phone","address","parent_name","account_number","account_branch","account_type","balance","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)).
		WithArgs(
			sqlmock.AnyArg(),
			customer.Name,
			customer.Phone,
			customer.Address,
			customer.ParentName,
			customer.AccountNumber,
			customer.AccountBranch,
			customer.AccountType,
			customer.Balance,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, customer)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, customer.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListCustomers(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := customer.NewCustomerRepository(db)
	ctx := context.Background()

	search := "test"
	limit := 10
	offset := 0

	mockCustomers := []model.Customer{
		{
			ID:            uuid.NewV4(),
			AccountNumber: "123456789",
			Name:          "Test Customer",
		},
	}

	t.Run("should return list of customers", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT count(*) FROM "customers" WHERE account_number ILIKE $1 OR name ILIKE $2`)).
			WithArgs("%"+search+"%", "%"+search+"%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(mockCustomers)))

		selectQuery := `SELECT * FROM "customers" WHERE account_number ILIKE $1 OR name ILIKE $2 LIMIT $3`
		selectArgs := []driver.Value{"%" + search + "%", "%" + search + "%", limit}
		if offset > 0 {
			selectQuery += " OFFSET $4"
			selectArgs = append(selectArgs, offset)
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
			WithArgs(selectArgs...).
			WillReturnRows(sqlmock.NewRows([]string{"id", "account_number", "name"}).
				AddRow(mockCustomers[0].ID, mockCustomers[0].AccountNumber, mockCustomers[0].Name))

		customers, total, err := repo.List(ctx, limit, offset, search)

		assert.NoError(t, err)
		assert.Equal(t, int64(len(mockCustomers)), total)
		assert.Len(t, customers, len(mockCustomers))

		if len(customers) > 0 {
			assert.Equal(t, mockCustomers[0].AccountNumber, customers[0].AccountNumber)
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when count fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT count(*) FROM "customers" WHERE account_number ILIKE $1 OR name ILIKE $2`)).
			WithArgs("%"+search+"%", "%"+search+"%").
			WillReturnError(assert.AnError)

		customers, total, err := repo.List(ctx, limit, offset, search)

		assert.Error(t, err)
		assert.Nil(t, customers)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT count(*) FROM "customers" WHERE account_number ILIKE $1 OR name ILIKE $2`)).
			WithArgs("%"+search+"%", "%"+search+"%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(mockCustomers)))

		selectQuery := `SELECT * FROM "customers" WHERE account_number ILIKE $1 OR name ILIKE $2 LIMIT $3`
		selectArgs := []driver.Value{"%" + search + "%", "%" + search + "%", limit}
		if offset > 0 {
			selectQuery += " OFFSET $4"
			selectArgs = append(selectArgs, offset)
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
			WithArgs(selectArgs...).
			WillReturnError(assert.AnError)

		customers, total, err := repo.List(ctx, limit, offset, search)

		assert.Error(t, err)
		assert.Nil(t, customers)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDetailCustomer(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := customer.NewCustomerRepository(db)
	ctx := context.Background()

	customerID := uuid.NewV4()

	expectedCustomer := model.Customer{
		ID:            customerID,
		Name:          "John Doe",
		AccountNumber: "1234567890",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Deposits: []model.Deposit{
			{
				ID:           uuid.NewV4(),
				CustomerID:   customerID,
				Amount:       1000.0,
				InterestRate: 3.5,
				TermMonths:   12,
				StartDate:    time.Now(),
				MaturityDate: time.Now().AddDate(1, 0, 0),
				IsWithdrawn:  false,
				Note:         "First deposit",
			},
		},
		Pockets: []model.Pocket{
			{
				ID:         uuid.NewV4(),
				CustomerID: customerID,
				Name:       "Travel",
				Balance:    500.0,
				IsActive:   true,
			},
		},
	}

	t.Run("should return detail of customer with deposits and pockets", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
			WithArgs(customerID, 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "account_number", "created_at", "updated_at",
			}).AddRow(expectedCustomer.ID, expectedCustomer.Name, expectedCustomer.AccountNumber, expectedCustomer.CreatedAt, expectedCustomer.UpdatedAt))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "deposits" WHERE "deposits"."customer_id" = $1`)).
			WithArgs(customerID).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "customer_id", "amount", "interest_rate", "term_months", "start_date", "maturity_date", "is_withdrawn", "note", "created_at", "updated_at",
			}).AddRow(
				expectedCustomer.Deposits[0].ID,
				customerID,
				expectedCustomer.Deposits[0].Amount,
				expectedCustomer.Deposits[0].InterestRate,
				expectedCustomer.Deposits[0].TermMonths,
				expectedCustomer.Deposits[0].StartDate,
				expectedCustomer.Deposits[0].MaturityDate,
				expectedCustomer.Deposits[0].IsWithdrawn,
				expectedCustomer.Deposits[0].Note,
				expectedCustomer.Deposits[0].CreatedAt,
				expectedCustomer.Deposits[0].UpdatedAt,
			))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pockets" WHERE "pockets"."customer_id" = $1`)).
			WithArgs(customerID).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "customer_id", "name", "balance", "target_amount", "target_date", "is_active", "created_at", "updated_at",
			}).AddRow(
				expectedCustomer.Pockets[0].ID,
				customerID,
				expectedCustomer.Pockets[0].Name,
				expectedCustomer.Pockets[0].Balance,
				nil,
				nil,
				expectedCustomer.Pockets[0].IsActive,
				expectedCustomer.Pockets[0].CreatedAt,
				expectedCustomer.Pockets[0].UpdatedAt,
			))

		result, err := repo.Detail(ctx, customerID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCustomer.Name, result.Name)
		assert.Equal(t, expectedCustomer.Deposits[0].Amount, result.Deposits[0].Amount)
		assert.Equal(t, expectedCustomer.Pockets[0].Name, result.Pockets[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error if customer not found", func(t *testing.T) {
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		repo := customer.NewCustomerRepository(db)
		ctx := context.Background()
		nonExistentID := uuid.NewV4()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
			WithArgs(nonExistentID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.Detail(ctx, nonExistentID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := customer.NewCustomerRepository(gormDB)

	id := "123"
	input := &model.Customer{
		Name:          "Updated Name",
		AccountNumber: "999999999",
		Balance:       2000.0,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "customers"`).
		WithArgs(input.Name, input.AccountNumber, input.Balance, AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	trx := gormDB.Begin()
	ctx := database.NewContext(context.Background(), trx)

	err = repo.Update(ctx, id, input)
	require.NoError(t, err)
	require.NoError(t, trx.Commit().Error)
}

func TestDeleteCustomer(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := customer.NewCustomerRepository(db)
	ctx := context.Background()

	customerID := uuid.NewV4()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "customers" WHERE id = $1`)).
		WithArgs(customerID.String()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, customerID.String())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := customer.NewCustomerRepository(gormDB)

	customerID := "123"
	amount := 500.0

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "customers"`).
		WithArgs(amount, AnyTime{}, customerID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	trx := gormDB.Begin()
	ctx := database.NewContext(context.Background(), trx)

	err = repo.AddBalance(ctx, customerID, amount)
	require.NoError(t, err)
	require.NoError(t, trx.Commit().Error)
}

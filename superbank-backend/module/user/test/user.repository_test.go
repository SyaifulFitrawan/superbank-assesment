package user

import (
	"bank-backend/model"
	"bank-backend/module/user"
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
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

func TestCreateUser(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := user.NewUserRepository(db)
	ctx := context.Background()

	u := &model.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpassword",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "users" ("id","email","username","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6)`)).
		WithArgs(sqlmock.AnyArg(), u.Email, u.Username, u.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, u)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, u.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := user.NewUserRepository(db)
	ctx := context.Background()

	search := "test"
	limit := 10
	offset := 0

	mockUsers := []model.User{
		{
			ID:       uuid.NewV4(),
			Email:    "test@example.com",
			Username: "testuser",
			Password: "hashedpassword",
		},
	}
	t.Run("should return list of users", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT count(*) FROM "users" WHERE email ILIKE $1 OR username ILIKE $2`)).
			WithArgs("%"+search+"%", "%"+search+"%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(mockUsers)))

		selectQuery := `SELECT * FROM "users" WHERE email ILIKE $1 OR username ILIKE $2 LIMIT $3`
		selectArgs := []driver.Value{"%" + search + "%", "%" + search + "%", limit}
		if offset > 0 {
			selectQuery += " OFFSET $4"
			selectArgs = append(selectArgs, offset)
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
			WithArgs(selectArgs...).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password"}).
				AddRow(mockUsers[0].ID, mockUsers[0].Email, mockUsers[0].Username, mockUsers[0].Password))

		users, total, err := repo.List(ctx, limit, offset, search)

		assert.NoError(t, err)
		assert.Equal(t, int64(len(mockUsers)), total)
		assert.Len(t, users, len(mockUsers))

		if len(users) > 0 {
			assert.Equal(t, mockUsers[0].Email, users[0].Email)
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when count fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT count(*) FROM "users" WHERE email ILIKE $1 OR username ILIKE $2`)).
			WithArgs("%"+search+"%", "%"+search+"%").
			WillReturnError(assert.AnError)

		users, total, err := repo.List(ctx, limit, offset, search)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when find fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT count(*) FROM "users" WHERE email ILIKE $1 OR username ILIKE $2`)).
			WithArgs("%"+search+"%", "%"+search+"%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(mockUsers)))

		selectQuery := `SELECT * FROM "users" WHERE email ILIKE $1 OR username ILIKE $2 LIMIT $3`
		selectArgs := []driver.Value{"%" + search + "%", "%" + search + "%", limit}
		if offset > 0 {
			selectQuery += " OFFSET $4"
			selectArgs = append(selectArgs, offset)
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
			WithArgs(selectArgs...).
			WillReturnError(assert.AnError)

		users, total, err := repo.List(ctx, limit, offset, search)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDetailUser(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := user.NewUserRepository(db)
	ctx := context.Background()

	userID := uuid.NewV4()
	expectedUser := model.User{
		ID:       userID,
		Email:    "detail@example.com",
		Username: "detailuser",
		Password: "secret",
	}

	t.Run("should return detail of users", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(userID, 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "email", "username", "password", "created_at", "updated_at",
			}).AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Username, expectedUser.Password, nil, nil))

		result, err := repo.Detail(ctx, userID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		repo := user.NewUserRepository(db)
		ctx := context.Background()

		nonExistentID := uuid.NewV4()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(nonExistentID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.Detail(ctx, nonExistentID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteUser(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := user.NewUserRepository(db)
	ctx := context.Background()

	userID := uuid.NewV4()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE id = $1`)).
		WithArgs(userID.String()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, userID.String())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByEmail(t *testing.T) {
	t.Run("should return user when email is found", func(t *testing.T) {
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		repo := user.NewUserRepository(db)
		ctx := context.Background()

		email := "found@example.com"
		expectedUser := model.User{
			ID:       uuid.NewV4(),
			Email:    email,
			Username: "founduser",
			Password: "secretpassword",
		}

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "email", "username", "password", "created_at", "updated_at",
			}).AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Username, expectedUser.Password, nil, nil))

		result, err := repo.FindByEmail(ctx, email)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.Username, result.Username)
		assert.Equal(t, expectedUser.Password, result.Password)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error if user not found by email", func(t *testing.T) {
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		repo := user.NewUserRepository(db)
		ctx := context.Background()

		email := "nonexistent@example.com"

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.FindByEmail(ctx, email)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

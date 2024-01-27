package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormUserRepository(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	repo := NewGormUserRepository(db)

	t.Run("FindByID", func(t *testing.T) {
		uid := uuid.New()

		rows := sqlmock.NewRows([]string{"uid", "email"}).
			AddRow(uid, "test@example.com")

		mock.ExpectQuery("SELECT").WithArgs(uid).WillReturnRows(rows)

		user, err := repo.FindByID(context.Background(), uid)

		assert.NoError(t, err)
		assert.Equal(t, uid, user.UID)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("FindByEmail", func(t *testing.T) {
		email := "test@example.com"

		rows := sqlmock.NewRows([]string{"uid", "email"}).
			AddRow(uuid.New(), email)

		mock.ExpectQuery("SELECT").WithArgs(email).WillReturnRows(rows)

		user, err := repo.FindByEmail(context.Background(), email)

		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
	})
}

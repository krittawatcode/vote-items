package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormVoteSessionRepository(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	repo := NewGormVoteSessionRepository(db)

	t.Run("GetOpenVoteSession", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "is_open"}).
			AddRow(1, true)

		mock.ExpectQuery("SELECT").WithArgs(true).WillReturnRows(rows)

		voteSession, err := repo.GetOpenVoteSession()

		assert.NoError(t, err)
		assert.Equal(t, uint(1), voteSession.ID)
		assert.Equal(t, true, voteSession.IsOpen)
	})

	t.Run("CreateVoteSession", func(t *testing.T) {
		id := uint(1)

		mock.ExpectExec("INSERT INTO").WithArgs(id, true).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateVoteSession(id)

		assert.NoError(t, err)
	})

	// TODO: CloseVoteSession

	t.Run("GetVoteSessionByID", func(t *testing.T) {
		id := uint(1)

		rows := sqlmock.NewRows([]string{"id", "is_open"}).
			AddRow(id, true)

		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(rows)

		voteSession, err := repo.GetVoteSessionByID(context.Background(), id)

		assert.NoError(t, err)
		assert.Equal(t, id, voteSession.ID)
		assert.Equal(t, true, voteSession.IsOpen)
	})
}

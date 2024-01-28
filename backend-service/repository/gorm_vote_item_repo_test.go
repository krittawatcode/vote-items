package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormVoteItemRepository(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	repo := NewGormVoteItemRepository(db)

	t.Run("FetchActive", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "vote_count", "session_id", "is_active"}).
			AddRow("3fa85f64-5717-4562-b3fc-2c963f66afa6", "Item 1", "Description 1", 10, 1, true).
			AddRow("3fa85f64-5717-4562-b3fc-2c963f66afa7", "Item 2", "Description 2", 20, 2, true)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		voteItems, err := repo.FetchActive(context.Background())

		assert.NoError(t, err)
		assert.Len(t, *voteItems, 2)
		assert.Equal(t, "Item 1", (*voteItems)[0].Name)
		assert.Equal(t, "Description 1", (*voteItems)[0].Description)
		assert.Equal(t, 10, (*voteItems)[0].VoteCount)
		assert.Equal(t, uint(1), (*voteItems)[0].SessionID)
		assert.Equal(t, true, (*voteItems)[0].IsActive)
		assert.Equal(t, "Item 2", (*voteItems)[1].Name)
		assert.Equal(t, "Description 2", (*voteItems)[1].Description)
		assert.Equal(t, 20, (*voteItems)[1].VoteCount)
		assert.Equal(t, uint(2), (*voteItems)[1].SessionID)
		assert.Equal(t, true, (*voteItems)[1].IsActive)
	})
}

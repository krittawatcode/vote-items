package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormVoteRepository(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	repo := NewGormVoteRepository(db)

	t.Run("GetVoteResultsBySession", func(t *testing.T) {
		sessionID := uint(1)

		rows := sqlmock.NewRows([]string{"vote_item_id", "vote_item_name", "vote_count"}).
			AddRow(uuid.New(), "Item 1", 10).
			AddRow(uuid.New(), "Item 2", 5)

		mock.ExpectQuery("SELECT").WithArgs(sessionID).WillReturnRows(rows)

		results, err := repo.GetVoteResultsBySession(sessionID)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Item 1", results[0].VoteItemName)
		assert.Equal(t, uint(10), results[0].VoteCount)
		assert.Equal(t, "Item 2", results[1].VoteItemName)
		assert.Equal(t, uint(5), results[1].VoteCount)
	})
}

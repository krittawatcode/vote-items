package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormVoteRepository_Create(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	repo := NewGormVoteRepository(db)

	userId := uuid.New()
	itemId := uuid.New()

	t.Run("No open vote session", func(t *testing.T) {
		vote := &domain.Vote{
			UserID:     userId,
			VoteItemID: itemId,
		}

		// Mock the vote session query
		mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)

		err := repo.Create(context.Background(), vote)

		assert.Error(t, err)
	})

	t.Run("User has already voted", func(t *testing.T) {
		vote := &domain.Vote{
			UserID:     userId,
			VoteItemID: itemId,
		}

		// Mock the vote session query
		mock.ExpectQuery("SELECT").WillReturnRows(
			sqlmock.NewRows([]string{"id", "is_open"}).AddRow(1, true),
		)

		// Mock the existing vote query
		mock.ExpectQuery("SELECT").WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "session_id", "vote_item_id"}).AddRow(1, vote.UserID, 1, vote.VoteItemID),
		)

		err := repo.Create(context.Background(), vote)

		assert.Error(t, err)
	})
}

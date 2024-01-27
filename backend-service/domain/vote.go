package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VoteSession struct {
	ID        uint `db:"id" json:"id"`
	IsOpen    bool `gorm:"type:boolean;not null;default:true" json:"is_open"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	BaseModel
}

type VoteSessionUseCase interface {
	GetOpenVoteSession() (*VoteSession, error)
	OpenVoteSession(id uint) error
	CloseVoteSession(id uint) error
	GetVoteSessionByID(ctx context.Context, id uint) (*VoteSession, error)
}
type VoteSessionRepository interface {
	GetOpenVoteSession() (*VoteSession, error)
	CreateVoteSession(id uint) error
	CloseVoteSession(id uint) error
	GetVoteSessionByID(ctx context.Context, id uint) (*VoteSession, error)
}

type VoteItem struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" binding:"required" json:"name"`
	Description string    `gorm:"type:text" binding:"required" json:"description"`
	VoteCount   int       `gorm:"type:int;default:0" json:"vote_count"`
	SessionID   uint      `gorm:"not null" json:"session_id"`
	IsActive    bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	BaseModel
}

// UserUseCase defines methods the handler layer expects
// any service it interacts with to implement
type VoteItemUseCase interface {
	FetchActive(ctx context.Context) (*[]VoteItem, error)
	Create(ctx context.Context, v *VoteItem) error
	Update(ctx context.Context, v *VoteItem) error
	Delete(ctx context.Context, vid uuid.UUID) error
	ClearVoteItem(ctx context.Context) error
}

// VoteItemRepository defines methods it expects a repository
// it interacts with to implement
type VoteItemRepository interface {
	FetchActive(ctx context.Context) (*[]VoteItem, error)
	Create(ctx context.Context, v *VoteItem) error
	Update(ctx context.Context, v *VoteItem) error
	SetActiveVoteItem(ctx context.Context, v *VoteItem, isActive bool) error
	ClearVoteItem(ctx context.Context) error
}

type Vote struct {
	BaseModel
	ID         uuid.UUID `db:"id" json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID     uuid.UUID `gorm:"not null" json:"user_id"`
	VoteItemID uuid.UUID `gorm:"type:uuid;not null" json:"vote_item_id"`
	SessionID  uint      `gorm:"not null" json:"session_id"`
}

type VoteResult struct {
	VoteItemID   uuid.UUID `json:"vote_item_id"`
	VoteItemName string    `json:"vote_item_name"`
	VoteCount    uint      `json:"vote_count" gorm:"column:vote_count"`
}

type VoteUseCase interface {
	Create(ctx context.Context, v *Vote) error
	GetVoteResultsBySession(sessionID uint) ([]VoteResult, error)
}

type VoteRepository interface {
	Create(ctx context.Context, v *Vote) error
	GetVoteResultsBySession(sessionID uint) ([]VoteResult, error)
}

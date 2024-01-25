package domain

import (
	"context"

	"github.com/google/uuid"
)

// User defines domain model and its json and db representations
type User struct {
	UID      uuid.UUID `db:"uid" json:"uid" gorm:"type:uuid;default:uuid_generate_v4()"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"-"` // never return password
}

// TableName use to specific table
func (u *User) TableName() string {
	return "users"
}

// UserUseCase defines methods the handler layer expects
// any service it interacts with to implement
type UserUseCase interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	SignUp(ctx context.Context, u *User) error
	SignIn(ctx context.Context, u *User) error
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
}

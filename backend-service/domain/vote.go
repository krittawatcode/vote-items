package domain

import "github.com/google/uuid"

type VoteItem struct {
	UID      uuid.UUID `db:"uid" json:"uid" gorm:"type:uuid;default:uuid_generate_v4()"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"-"` // never return password
}

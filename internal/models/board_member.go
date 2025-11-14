package models

import (
	"time"
)

type BoardMember struct {
	ID         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	BoardId    int64     `json:"board_id"`
	MemberRole string    `json:"member_role"`
	CreatedAt  time.Time `json:"created_at"` // TODO: Добавить CreatedAt
	UpdatedAt  time.Time `json:"updated_at"` // TODO: Добавить UpdateAt
}

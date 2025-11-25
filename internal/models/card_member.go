package models

import (
	"time"
)

type CardMember struct {
	ID            int64     `json:"id"`
	CardID        int64     `json:"card_id"`
	BoardMemberID int64     `json:"board_member_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

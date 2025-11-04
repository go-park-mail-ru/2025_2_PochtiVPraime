package models

import "time"

type Comment struct {
	ID                 int64     `json:"id"`
	CardId             int64     `json:"card_id"`
	BoardMemberOwnerId int64     `json:"board_member_owner_id"`
	Content            string    `json:"content"`
	CreatedAt          time.Time `json:"created_at"` //?
	UpdatedAt          time.Time `json:"updated_at"` //?
}
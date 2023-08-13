package entities

import (
	"time"
	. "zhasa2.0/user/entities"
)

type Comment struct {
	CommentId   int32
	Message     string
	User        User
	CreatedDate time.Time
	Id          int32
}

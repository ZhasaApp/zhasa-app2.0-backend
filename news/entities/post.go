package entities

import (
	"time"
	. "zhasa2.0/user/entities"
)

type Post struct {
	Id            int32
	Images        []string
	LikesCount    int32
	CommentsCount int32
	Title         string
	Body          string
	IsLiked       bool
	LikesByOwner  int32
	Author        User
	CreatedDate   time.Time
}

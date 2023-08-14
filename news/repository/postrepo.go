package repository

import (
	"context"
	"database/sql"
	"fmt"
	. "zhasa2.0/base"
	. "zhasa2.0/db/hand-made"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/news/entities"
	. "zhasa2.0/user/entities"
)

type PostRepository interface {
	CreatePost(postTitle, postBody string, authorId int32, imageUrls []string) error
	CreateComment(userId int32, postId int32, message string) error
	GetPostComments(postId int32, pagination Pagination) ([]Comment, error)
	GetPosts(userId int32, pagination Pagination) ([]Post, error)
	AddLike(userId int32, postId int32) error
	DeleteLike(userId int32, postId int32) error
	IsUserLikedPost(userId int32, postId int32) (bool, error)
	DeletePost(postId int32) error
	DeleteComment(commentId int32) error
}

type DBPostRepository struct {
	ctx     context.Context
	querier generated.Querier
	customQ CustomQuerier
}

func (db DBPostRepository) DeleteComment(commentId int32) error {
	err := db.querier.DeleteComment(db.ctx, commentId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (db DBPostRepository) GetPostComments(postId int32, pagination Pagination) ([]Comment, error) {
	rows, err := db.querier.GetCommentsAndAuthorsByPostId(db.ctx, generated.GetCommentsAndAuthorsByPostIdParams{
		PostID: postId,
		Limit:  pagination.PageSize,
		Offset: pagination.Page,
	})

	comments := make([]Comment, 0)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, row := range rows {
		comments = append(comments, Comment{
			CommentId: row.CommentID,
			Message:   row.Body,
			User: User{
				Id:        row.UserID,
				Phone:     "",
				Avatar:    row.AvatarUrl,
				FirstName: Name(row.FirstName),
				LastName:  Name(row.LastName),
			},
			CreatedDate: row.CreatedAt,
			Id:          row.CommentID,
		})
	}

	return comments, nil
}

func (db DBPostRepository) CreateComment(userId int32, postId int32, message string) error {
	_, err := db.querier.CreateComment(db.ctx, generated.CreateCommentParams{
		Body:   message,
		UserID: userId,
		PostID: postId,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (db DBPostRepository) DeletePost(postId int32) error {
	err := db.querier.DeletePost(db.ctx, postId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (db DBPostRepository) AddLike(userId int32, postId int32) error {
	_, err := db.querier.AddLike(db.ctx, generated.AddLikeParams{
		UserID: userId,
		PostID: postId,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (db DBPostRepository) DeleteLike(userId int32, postId int32) error {
	err := db.querier.DeleteLike(db.ctx, generated.DeleteLikeParams{
		UserID: userId,
		PostID: postId,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (db DBPostRepository) IsUserLikedPost(userId int32, postId int32) (bool, error) {
	_, err := db.querier.GetUserPostLike(db.ctx, generated.GetUserPostLikeParams{
		UserID: userId,
		PostID: postId,
	})

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (db DBPostRepository) GetPosts(userId int32, pagination Pagination) ([]Post, error) {
	rows, err := db.customQ.GetPostsAndPostAuthors(db.ctx, GetPostsAndPostAuthorsParams{
		UserID: userId,
		Limit:  pagination.PageSize,
		Offset: pagination.Page,
	})

	posts := make([]Post, 0)

	if err == sql.ErrNoRows {
		return posts, nil
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, row := range rows {
		posts = append(posts, Post{
			Id:            row.ID,
			Images:        row.ImageUrls,
			LikesCount:    int32(row.LikesCount),
			CommentsCount: int32(row.CommentsCount),
			Title:         row.Title,
			Body:          row.Body,
			IsLiked:       row.IsLiked,
			Author: User{
				Id:        row.UserID,
				Avatar:    row.AvatarUrl,
				FirstName: Name(row.FirstName),
				LastName:  Name(row.LastName),
			},
			CreatedDate: row.CreatedAt,
		})
	}

	return posts, nil
}

func (db DBPostRepository) CreatePost(postTitle, postBody string, authorId int32, imageUrls []string) error {
	post, err := db.querier.CreatePost(db.ctx, generated.CreatePostParams{
		Title:  postTitle,
		Body:   postBody,
		UserID: authorId,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, url := range imageUrls {
		err := db.querier.CreatePostImages(db.ctx, generated.CreatePostImagesParams{
			ImageUrl: url,
			PostID:   post.ID,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func NewPostRepository(ctx context.Context, querier generated.Querier) PostRepository {
	return DBPostRepository{
		ctx:     ctx,
		querier: querier,
	}
}

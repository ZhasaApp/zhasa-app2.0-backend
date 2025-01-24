package hand_made

import (
	"context"
	"github.com/lib/pq"
	"time"
)

const getPostsAndPostAuthors = `-- name: GetPostsAndPostAuthors :many
SELECT p.id, p.title, p.body, p.user_id, p.created_at,
       EXISTS(SELECT user_id, post_id FROM likes l WHERE l.post_id = p.id AND l.user_id = $1) AS is_liked,
       COALESCE(lc.likes_count, 0)                                             AS likes_count,
       COALESCE(cc.comments_count, 0)                                          AS comments_count,
       COALESCE(
               (SELECT ARRAY_AGG(p_i.image_url)
                FROM post_images p_i
                WHERE p_i.post_id = p.id),
               ARRAY[] ::text[]
           )                                                                   AS image_urls,
       u.id                                                                    AS user_id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       COALESCE(
               (SELECT COUNT(*)
                FROM likes l
                         JOIN user_roles ur ON l.user_id = ur.user_id
                         JOIN roles r ON ur.role_id = r.id
                WHERE l.post_id = p.id
                  AND r.key = 'owner'),
               0
       )                                                                        AS likes_by_owner
FROM (SELECT id, title, body, user_id, created_at FROM posts ORDER BY created_at DESC LIMIT $2 OFFSET $3) p
         LEFT JOIN
         (SELECT post_id, COUNT(*) AS likes_count FROM likes GROUP BY post_id) lc ON lc.post_id = p.id
         LEFT JOIN
     (SELECT post_id, COUNT(*) AS comments_count FROM comments GROUP BY post_id) cc ON cc.post_id = p.id
         JOIN
     user_avatar_view u ON p.user_id = u.id
ORDER BY p.created_at DESC
`

type GetPostsAndPostAuthorsParams struct {
	UserID int32 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetPostsAndPostAuthorsRow struct {
	ID            int32          `json:"id"`
	Title         string         `json:"title"`
	Body          string         `json:"body"`
	UserID        int32          `json:"user_id"`
	CreatedAt     time.Time      `json:"created_at"`
	IsLiked       bool           `json:"is_liked"`
	LikesCount    int64          `json:"likes_count"`
	CommentsCount int64          `json:"comments_count"`
	ImageUrls     pq.StringArray `json:"image_urls"`
	UserID_2      int32          `json:"user_id_2"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	AvatarUrl     string         `json:"avatar_url"`
	LikesByOwner  int32          `json:"likes_by_owner"`
}

func (q DBCustomQuerier) GetPostsAndPostAuthors(ctx context.Context, arg GetPostsAndPostAuthorsParams) ([]GetPostsAndPostAuthorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsAndPostAuthors, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsAndPostAuthorsRow
	for rows.Next() {
		var i GetPostsAndPostAuthorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Body,
			&i.UserID,
			&i.CreatedAt,
			&i.IsLiked,
			&i.LikesCount,
			&i.CommentsCount,
			&i.ImageUrls,
			&i.UserID_2,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.LikesByOwner,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostsAndPostAuthorsCount = `-- name: GetPostsAndPostAuthorsCount :one
select count(*) from posts
`

func (q DBCustomQuerier) GetPostsAndPostAuthorsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPostsAndPostAuthorsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

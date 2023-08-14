-- name: GetPostById :one
SELECT *
FROM posts
WHERE id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT *
FROM posts
ORDER BY created_at;

-- name: CreatePost :one
INSERT INTO posts (title, body, user_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: CreatePostImages :exec
INSERT INTO post_images (image_url, post_id)
VALUES ($1, $2);

-- name: DeletePost :exec
DELETE
FROM posts
WHERE id = $1;

-- name: GetPostsAndPostAuthors :many
SELECT p.*,
       EXISTS(SELECT * FROM likes l WHERE l.post_id = p.id AND l.user_id = $1) AS is_liked,
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
       u.avatar_url
FROM (SELECT * FROM posts ORDER BY created_at DESC LIMIT $2 OFFSET $3) p
         LEFT JOIN
         (SELECT post_id, COUNT(*) AS likes_count FROM likes GROUP BY post_id) lc ON lc.post_id = p.id
         LEFT JOIN
     (SELECT post_id, COUNT(*) AS comments_count FROM comments GROUP BY post_id) cc ON cc.post_id = p.id
         JOIN
     user_avatar_view u ON p.user_id = u.id
ORDER BY p.created_at DESC;


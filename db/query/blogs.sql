-- name: CreateBlog :one
INSERT INTO blogs (
  username,
  title,
  content
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetBlog :one
SELECT * FROM blogs
WHERE id = $1 LIMIT 1;
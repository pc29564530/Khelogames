-- name: getSearchBar :many
SELECT full_name, username, communities FROM search_bar
WHERE (full_name ILIKE '%' || $1 || '%') | (username ILIKE '%' || $1 || '%') | (communities ILIKE '%' || $1 || '%');

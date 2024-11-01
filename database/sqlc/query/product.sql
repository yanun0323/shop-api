-- name: ListProducts :many
SELECT * FROM `product` 
WHERE `category_id` = ? 
ORDER BY `rank` DESC;
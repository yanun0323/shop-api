-- name: GetUserCategoryPreferenceByUserID :one
SELECT `category_id` FROM `user_category_preference` WHERE `user_id` = ? LIMIT 1;
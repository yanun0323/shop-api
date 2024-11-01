-- name: CountToken :one
SELECT COUNT(*) FROM `token` 
WHERE `user_id` = ? AND `device_id` = ? AND `expired_at` > ?;

-- name: GetToken :one
SELECT * FROM `token` 
WHERE `user_id` = ? AND `device_id` = ? AND `expired_at` > ?;

-- name: CreateToken :exec
INSERT INTO `token` (`user_id`, `device_id`, `refresh_token`, `expired_at`) 
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE 
    `refresh_token` = VALUES(`refresh_token`),
    `expired_at` = VALUES(`expired_at`),
    `updated_at` = UNIX_TIMESTAMP();

-- name: DeleteToken :exec
DELETE FROM `token` WHERE `user_id` = ? AND `device_id` = ?;
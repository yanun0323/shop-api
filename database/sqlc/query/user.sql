-- name: CountUser :one
SELECT COUNT(*) FROM `user` WHERE `email` = ? LIMIT 1;

-- name: CreateUser :execlastid
INSERT INTO `user` (`name`, `email`, `password`) 
VALUES (?, ?, ?);

-- name: GetUserByEmail :one
SELECT * FROM `user` WHERE `email` = ? LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM `user` WHERE `id` = ?;
-- +goose Up
CREATE TABLE IF NOT EXISTS `token` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `device_id` varchar(64) NOT NULL,
  `refresh_token` varchar(255) NOT NULL,
  `created_at` bigint NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  `updated_at` bigint NULL,
  `expired_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `index_user_id_device_id_expired_at` (`user_id`, `device_id`, `expired_at`),
  UNIQUE KEY `unique_user_id_device_id` (`user_id`, `device_id`)
);

-- +goose Down 
DROP TABLE IF EXISTS `token`;
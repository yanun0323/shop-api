-- +goose Up
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` char(64) NOT NULL COMMENT 'SHA-256 加密後的密碼',
  `created_at` bigint NOT NULL,
  `updated_at` bigint NULL,
  `deleted_at` bigint NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_email` (`email`)
);

-- +goose Down
DROP TABLE IF EXISTS `user`;
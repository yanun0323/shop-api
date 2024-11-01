-- +goose Up
CREATE TABLE IF NOT EXISTS `user_category_preference` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `category_id` bigint NOT NULL,
  `created_at` bigint NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  `updated_at` bigint NULL,
  PRIMARY KEY (`id`),
  INDEX `index_user_id_category_id` (`user_id`, `category_id`)
);

-- +goose Down
DROP TABLE IF EXISTS `user_category_preference`;
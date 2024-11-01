-- +goose Up
CREATE TABLE IF NOT EXISTS `product` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `category_id` bigint NOT NULL,
  `price` decimal(10, 2) NOT NULL,
  `rank` int NOT NULL COMMENT '推薦排序',
  `created_at` bigint NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  `updated_at` bigint NULL,
  PRIMARY KEY (`id`),
  INDEX `index_category_id_rank` (`category_id`, `rank`) COMMENT '使用者推薦使用'
);

INSERT INTO product (`name`, `description`, `category_id`, `price`, `rank`, `created_at`, `updated_at`) VALUES
('iPhone 15', '最新款智慧型手機', 1, 799.99, 1, UNIX_TIMESTAMP(), NULL),
('MacBook Air', '輕薄筆電', 1, 999.99, 2, UNIX_TIMESTAMP(), NULL),
('AirPods Pro', '無線耳機', 1, 249.99, 3, UNIX_TIMESTAMP(), NULL),
('Nike Air Max', '運動鞋', 2, 129.99, 4, UNIX_TIMESTAMP(), NULL),
('Adidas Ultra Boost', '跑步鞋', 2, 159.99, 5, UNIX_TIMESTAMP(), NULL),
('Levis 501', '牛仔褲', 2, 59.99, 6, UNIX_TIMESTAMP(), NULL),
('Samsung TV', '4K智慧電視', 3, 699.99, 7, UNIX_TIMESTAMP(), NULL),
('Sony PS5', '遊戲主機', 3, 499.99, 8, UNIX_TIMESTAMP(), NULL),
('Nintendo Switch', '掌上遊戲機', 3, 299.99, 9, UNIX_TIMESTAMP(), NULL),
('iPad Pro', '平板電腦', 1, 899.99, 10, UNIX_TIMESTAMP(), NULL),
('Kindle', '電子書閱讀器', 1, 139.99, 11, UNIX_TIMESTAMP(), NULL),
('Under Armour 運動衫', '運動服飾', 2, 45.99, 12, UNIX_TIMESTAMP(), NULL),
('Bose QC45', '降噪耳機', 1, 329.99, 13, UNIX_TIMESTAMP(), NULL),
('Canon EOS R', '數位相機', 1, 1999.99, 14, UNIX_TIMESTAMP(), NULL),
('Nike 運動褲', '運動服飾', 2, 49.99, 15, UNIX_TIMESTAMP(), NULL),
('Samsung Galaxy S23', '智慧型手機', 1, 699.99, 16, UNIX_TIMESTAMP(), NULL),
('Apple Watch', '智慧手錶', 1, 399.99, 17, UNIX_TIMESTAMP(), NULL),
('ASUS ROG', '電競筆電', 1, 1499.99, 18, UNIX_TIMESTAMP(), NULL),
('Dyson V15', '無線吸塵器', 3, 699.99, 19, UNIX_TIMESTAMP(), NULL),
('LG 冰箱', '智慧冰箱', 3, 1299.99, 20, UNIX_TIMESTAMP(), NULL);

-- +goose Down
DROP TABLE IF EXISTS `product`;
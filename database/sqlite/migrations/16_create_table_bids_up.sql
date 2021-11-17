CREATE TABLE `bids` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `provider_id` VARCHAR(255) NOT NULL,
  `booking_id` VARCHAR(255) NOT NULL,
  `amount` INTEGER NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
  -- FOREIGN KEY(`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
);
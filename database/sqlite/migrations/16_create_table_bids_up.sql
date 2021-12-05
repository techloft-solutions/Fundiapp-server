CREATE TABLE `bids` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `provider_id` VARCHAR(255) NOT NULL,
  `booking_id` VARCHAR(255) NOT NULL,
  `amount` INTEGER NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`provider_id`) REFERENCES `providers` (`provider_id`),
  FOREIGN KEY(`booking_id`) REFERENCES `bookings`(`booking_id`)
);
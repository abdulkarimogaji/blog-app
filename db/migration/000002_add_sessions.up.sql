CREATE TABLE `sessions` (
  `id` varchar(36) PRIMARY KEY,
  `user_id` int,
  `refresh_token` varchar(500),
  `client_ip` varchar(255),
  `user_agent` varchar(255),
  `is_blocked` boolean DEFAULT false,
  `created_at` timestamp,
  `expires_at` timestamp
);

ALTER TABLE `sessions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`); 
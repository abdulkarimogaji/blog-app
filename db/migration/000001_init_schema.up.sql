CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `first_name` varchar(255),
  `last_name` varchar(255),
  `email` varchar(255) UNIQUE,
  `password` text,
  `is_email_verified` boolean,
  `created_at` timestamp,
  `updated_at` timestamp
);

CREATE TABLE `profile` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int,
  `about` text,
  `date_of_birth` date,
  `photo` text,
  `city` varchar(255),
  `country` varchar(255),
  `settings` json,
  `socials` json,
  `created_at` timestamp,
  `updated_at` timestamp
);

CREATE TABLE `blogs` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `author_id` int,
  `title` text,
  `slug` text,
  `excerpt` text,
  `thumbnail` text,
  `body` longtext,
  `posted_at` timestamp,
  `created_at` timestamp,
  `updated_at` timestamp
);

CREATE TABLE `comments` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `blog_id` int,
  `user_id` int,
  `message` text,
  `thread` json,
  `created_at` timestamp,
  `updated_at` timestamp
);

CREATE TABLE `verify_email` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int,
  `email` varchar(255),
  `is_used` boolean,
  `secret_code` text,
  `created_at` timestamp,
  `expired_at` timestamp
);

ALTER TABLE `profile` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `blogs` ADD FOREIGN KEY (`author_id`) REFERENCES `users` (`id`);

ALTER TABLE `verify_email` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`blog_id`) REFERENCES `blogs` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

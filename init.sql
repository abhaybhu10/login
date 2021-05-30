CREATE DATABASE IF NOT EXISTS login;

USE login;
CREATE TABLE IF NOT EXISTS `users` (
  `id` VARCHAR(255),
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY(`id`)
);
CREATE TABLE IF NOT EXISTS `sessions` (
  `session_id` VARCHAR(255),
  `user_id` VARCHAR(255) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY(`session_id`)
)
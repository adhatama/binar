CREATE TABLE IF NOT EXISTS `product` (
	`id` TEXT,
	`name` TEXT,
    `price` INTEGER,
    `imageurl` TEXT,
	`created_at` DATETIME,
	`updated_at` DATETIME,
	PRIMARY KEY(`id`)
);

CREATE TABLE IF NOT EXISTS `user` (
	`id` TEXT,
	`name` TEXT,
    `email` TEXT,
    `password` TEXT,
	`created_at` DATETIME,
	`updated_at` DATETIME,
	PRIMARY KEY(`id`)
);

CREATE TABLE IF NOT EXISTS `user_auth` (
	`id` TEXT,
	`user_id` TEXT,
    `access_token` TEXT,
    `expired_at` DATETIME,
	`created_at` DATETIME,
	`updated_at` DATETIME,
	PRIMARY KEY(`id`),
	FOREIGN KEY(`user_id`) REFERENCES `user`(`id`) ON UPDATE CASCADE
);
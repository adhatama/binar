CREATE TABLE IF NOT EXISTS `product` (
	`id` TEXT,
	`name` TEXT,
    `price` INTEGER,
    `imageurl` TEXT,
	`created_at` DATETIME,
	`updated_at` DATETIME,
	PRIMARY KEY(`id`)
);
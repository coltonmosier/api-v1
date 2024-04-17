CREATE DATABASE `devices`;

CREATE TABLE `device_type` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(13) NOT NULL,
  `status` enum('active','inactive') NOT NULL DEFAULT 'active',
  PRIMARY KEY (`id`)
);

CREATE TABLE `manufacturer` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(10) NOT NULL,
  `status` enum('active','inactive') NOT NULL DEFAULT 'active',
  PRIMARY KEY (`id`)
);

CREATE TABLE `serial_numbers` (
  `auto_id` int NOT NULL AUTO_INCREMENT,
  `device_type_id` int NOT NULL,
  `manufacturer_id` int NOT NULL,
  `serial_number` varchar(68) NOT NULL,
  `status` enum('active','inactive') NOT NULL DEFAULT 'active',
  UNIQUE KEY `serial_number` (`serial_number`),
  KEY `device_type_id` (`device_type_id`),
  KEY `manufacturer_id` (`manufacturer_id`),
  CONSTRAINT `fk_to_device_type` FOREIGN KEY (`device_type_id`) REFERENCES `device_type` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_to_manufacturer` FOREIGN KEY (`manufacturer_id`) REFERENCES `manufacturer` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE SCHEMA `crypto` DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci;

CREATE TABLE `crypto`.`log_cross_exchange_tick` (
		  `id` INT NOT NULL AUTO_INCREMENT,
		  `symbol` VARCHAR(45) NULL,
		  `ask_exchange` INT NULL,
		  `ask_c_price` DOUBLE NULL,
		  `ask_s_price` DOUBLE NULL,
		  `ask_total_volume` DOUBLE NULL,
		  `bid_exchange` INT NULL,
		  `bid_c_price` DOUBLE NULL,
		  `bid_s_price` DOUBLE NULL,
		  `bid_total_volume` DOUBLE NULL,
		  `profit` DOUBLE NULL,
		  `min_total_volume` DOUBLE NULL,
		  `create_time` BIGINT NULL ,
		  PRIMARY KEY (`id`))
		COMMENT = 'save every x ms data';

CREATE INDEX `idx_log_cross_exchange_tick_profit`  ON `crypto`.`log_cross_exchange_tick` (profit) COMMENT '' ALGORITHM DEFAULT LOCK DEFAULT
CREATE INDEX `idx_log_cross_exchange_tick_create_time`  ON `crypto`.`log_cross_exchange_tick` (create_time) COMMENT '' ALGORITHM DEFAULT LOCK DEFAULT

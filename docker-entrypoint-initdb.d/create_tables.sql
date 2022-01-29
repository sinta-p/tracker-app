CREATE TABLE IF NOT EXISTS `tracker_db`.`stocks_tab` (
        `ticker` VARCHAR(8) NOT NULL,
        `company` VARCHAR(64) NOT NULL,
        `description` VARCHAR(1024),
        PRIMARY KEY (`ticker`),
        UNIQUE INDEX (`company`)
)
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB;

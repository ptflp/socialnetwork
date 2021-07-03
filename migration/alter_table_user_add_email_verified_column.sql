ALTER TABLE `infoblog`.`users`
    ADD COLUMN `email_verified` INT NULL DEFAULT 0 AFTER `password`;
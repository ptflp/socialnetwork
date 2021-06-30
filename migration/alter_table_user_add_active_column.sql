ALTER TABLE `infoblog`.`users`
    ADD COLUMN `active` INT NULL DEFAULT 0 AFTER `password`;
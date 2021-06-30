ALTER TABLE `infoblog`.`users`
    ADD COLUMN `password` VARCHAR(60) NULL AFTER `email`;
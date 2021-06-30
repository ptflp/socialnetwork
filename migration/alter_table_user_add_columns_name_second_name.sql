ALTER TABLE `infoblog`.`users`
    ADD COLUMN `name` VARCHAR(55) NULL DEFAULT NULL AFTER `active`,
    ADD COLUMN `second_name` VARCHAR(55) NULL DEFAULT NULL AFTER `name`;

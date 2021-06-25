CREATE TABLE `infoblog`.`users` (
                                        `id` INT NOT NULL AUTO_INCREMENT,
                                        `phone` VARCHAR(34) NULL,
                                        `email` VARCHAR(89) NULL,
                                        PRIMARY KEY (`id`),
                                        UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
                                        UNIQUE INDEX `phone_UNIQUE` (`phone` ASC) VISIBLE,
                                        UNIQUE INDEX `email_UNIQUE` (`email` ASC) VISIBLE);
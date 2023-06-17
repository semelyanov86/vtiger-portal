-- CREATE TABLE "notifications" --------------------------------
CREATE TABLE `notifications`(
                                `id` Int( 255 ) AUTO_INCREMENT PRIMARY KEY NOT NULL,
                                `crmid` VarChar( 50 ) NOT NULL,
                                `module` VarChar( 50 ) NOT NULL,
                                `label` VarChar( 255 ) NOT NULL,
                                `description` Text NOT NULL,
                                `assigned_user_id` VarChar( 50 ) NOT NULL,
                                `account_id` VarChar( 50 ) NOT NULL,
                                `user_id` VarChar( 50 ) NOT NULL,
                                `is_read` TinyInt( 255 ) NOT NULL DEFAULT 0,
                                `parent_id` VarChar( 50 ) NOT NULL,
                                created_at TIMESTAMP NOT NULL,
                                updated_at TIMESTAMP NOT NULL,
                                CONSTRAINT `unique_id` UNIQUE( `id` ) )
    ENGINE = InnoDB;
-- -------------------------------------------------------------

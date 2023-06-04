-- CREATE FIELD "otp_enabled" ----------------------------------
ALTER TABLE `users` ADD COLUMN `otp_enabled` TinyInt NOT NULL DEFAULT 0;
-- -------------------------------------------------------------

-- CREATE FIELD "otp_verified" ---------------------------------
ALTER TABLE `users` ADD COLUMN `otp_verified` TinyInt( 255 ) NOT NULL DEFAULT 0;
-- -------------------------------------------------------------

-- CREATE FIELD "otp_secret" -----------------------------------
ALTER TABLE `users` ADD COLUMN `otp_secret` VarChar( 190 ) NOT NULL DEFAULT '';
-- -------------------------------------------------------------

-- CREATE FIELD "otp_auth_url" ---------------------------------
ALTER TABLE `users` ADD COLUMN `otp_auth_url` VarChar( 255 ) NOT NULL DEFAULT '';
-- -------------------------------------------------------------

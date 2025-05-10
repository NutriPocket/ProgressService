CREATE DATABASE IF NOT EXISTS mydb; 

USE mydb;

SET time_zone = '+00:00';

CREATE TABLE IF NOT EXISTS fixed_user_data (
    user_id VARCHAR(36) PRIMARY KEY,
    height SMALLINT UNSIGNED NOT NULL,
    birthday DATE NOT NULL,
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE TABLE IF NOT EXISTS anthropometric_data (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    weight DECIMAL(5,2) NOT NULL,
    muscle_mass DECIMAL(5,2),
    fat_mass DECIMAL(5,2),
    bone_mass DECIMAL(5,2),
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6)
);

-- Delete expired JWTs

DELIMITER //

CREATE EVENT IF NOT EXISTS delete_expired_blacklist
ON SCHEDULE EVERY 1 MINUTE
DO
BEGIN
    DELETE FROM jwt_blacklist WHERE expires_at < NOW();
END //

DELIMITER ;
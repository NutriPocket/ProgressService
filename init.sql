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

CREATE TABLE IF NOT EXISTS objective (
    user_id VARCHAR(36) PRIMARY KEY,
    weight DECIMAL(5,2) NOT NULL,
    muscle_mass DECIMAL(5,2),
    fat_mass DECIMAL(5,2),
    bone_mass DECIMAL(5,2),
    deadline DATE NOT NULL,
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
);
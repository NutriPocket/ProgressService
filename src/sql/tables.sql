CREATE TABLE IF NOT EXISTS fixed_user_data (
    user_id VARCHAR(36) PRIMARY KEY,
    height SMALLINT UNSIGNED NOT NULL,
    birthday DATE NOT NULL,
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE TABLE IF NOT EXISTS anthropometric_data (
    user_id VARCHAR(36) NOT NULL,
    weight DECIMAL(5,2) NOT NULL,
    muscle_mass DECIMAL(5,2),
    fat_mass DECIMAL(5,2),
    bone_mass DECIMAL(5,2),
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL,
    PRIMARY KEY (user_id, created_at)
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

CREATE TABLE IF NOT EXISTS user_routines (
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(64) NOT NULL,
    description VARCHAR(512),
    day VARCHAR(10) NOT NULL,
    start_hour SMALLINT NOT NULL,
    end_hour SMALLINT NOT NULL,
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (user_id, day, start_hour, end_hour)
);

CREATE TABLE IF NOT EXISTS exercise_by_day (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    exercise_name VARCHAR(64) NOT NULL,
    calories_burned DECIMAL(6,2) NOT NULL,
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6)
);
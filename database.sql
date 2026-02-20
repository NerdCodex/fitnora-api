CREATE SCHEMA fitnora;
USE fitnora;
drop table users;
CREATE TABLE IF NOT EXISTS users (
    user_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    user_email VARCHAR(255) NOT NULL UNIQUE,
    user_fullname VARCHAR(100) NOT NULL,
    user_dob DATE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,

    gender VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS data_backup (
    data_backup_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    user_dbfiles TEXT NOT NULL,
    user_images TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

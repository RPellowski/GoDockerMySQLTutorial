CREATE DATABASE IF NOT EXISTS LOTRdata;
CREATE TABLE IF NOT EXISTS LOTRdata.users(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50),
    password VARCHAR(120)
);

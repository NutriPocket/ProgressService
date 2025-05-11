CREATE DATABASE IF NOT EXISTS test; 

USE test;

SET time_zone = '+00:00';

SHOW VARIABLES LIKE 'datadir';

SOURCE ./tables.sql;
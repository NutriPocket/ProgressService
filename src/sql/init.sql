CREATE DATABASE IF NOT EXISTS mydb; 

USE mydb;

SET time_zone = '+00:00';

SOURCE /src/sql/tables.sql;
SOURCE /src/sql/default.sql;
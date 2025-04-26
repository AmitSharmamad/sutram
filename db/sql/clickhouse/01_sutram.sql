CREATE DATABASE IF NOT EXISTS sutram;

USE sutram;

CREATE USER IF NOT EXISTS test IDENTIFIED WITH plaintext_password BY 'test';

GRANT ALL ON sutram.* TO test;

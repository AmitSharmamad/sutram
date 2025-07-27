-- Drop tables in reverse order to handle foreign key dependencies
DROP TABLE IF EXISTS user_permission;
DROP TABLE IF EXISTS role_permission;
DROP TABLE IF EXISTS permission;
DROP TABLE IF EXISTS user_role;
DROP TABLE IF EXISTS role;
DROP TABLE IF EXISTS "user";

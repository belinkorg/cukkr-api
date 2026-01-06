DROP CONSTRAINT IF EXISTS fk_user_barbershops_barbershop;
DROP CONSTRAINT IF EXISTS fk_user_barbershops_user;
DROP INDEX IF EXISTS idx_user_barbershops_role;
DROP INDEX IF EXISTS idx_user_barbershops_barbershop_id;
DROP INDEX IF EXISTS idx_user_barbershops_user_id;
DROP TABLE IF EXISTS user_barbershops;
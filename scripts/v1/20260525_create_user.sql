-- =========================================
-- 1. Tạo user
-- =========================================
CREATE USER app_otp
WITH PASSWORD 'app_otp2026';

-- =========================================
-- 2. Tạo database
-- =========================================
CREATE DATABASE otp
OWNER app_otp;

-- =========================================
-- 3. Kết nối vào DB otp
-- =========================================
\c otp

-- =========================================
-- 4. Grant quyền schema public
-- =========================================
GRANT USAGE, CREATE ON SCHEMA public TO app_otp;

-- =========================================
-- 5. Grant toàn quyền trên các table hiện có
-- =========================================
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO app_otp;

-- =========================================
-- 6. Grant toàn quyền trên sequence
-- =========================================
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO app_otp;

-- =========================================
-- 7. Grant quyền mặc định cho table tạo sau này
-- =========================================
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT ALL PRIVILEGES ON TABLES TO app_otp;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT ALL PRIVILEGES ON SEQUENCES TO app_otp;
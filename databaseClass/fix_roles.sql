-- 修复Role表，确保包含所有必需的中文角色
SET NAMES utf8mb4;
SET CHARACTER_SET_CLIENT = utf8mb4;
SET CHARACTER_SET_CONNECTION = utf8mb4;
SET CHARACTER_SET_RESULTS = utf8mb4;

USE volunteer;

-- 禁用外键约束检查
SET FOREIGN_KEY_CHECKS = 0;

-- 删除现有角色数据
DELETE FROM Role;

-- 插入正确的中文角色名
INSERT INTO Role (role_id, role_name) VALUES 
(1, '普通用户'),
(2, '管理员');

-- 重新启用外键约束检查
SET FOREIGN_KEY_CHECKS = 1;

-- 验证插入结果
SELECT * FROM Role;

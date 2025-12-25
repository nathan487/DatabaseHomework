-- 初始化数据脚本
-- 用于volunteer数据库的基础数据初始化

SET NAMES utf8mb4;
SET CHARACTER_SET_CLIENT = utf8mb4;
SET CHARACTER_SET_CONNECTION = utf8mb4;
SET CHARACTER_SET_RESULTS = utf8mb4;

USE volunteer;

-- 插入角色数据
INSERT IGNORE INTO Role (role_id, role_name) VALUES 
(1, '普通用户'),
(2, '管理员');

-- 插入部门数据
INSERT IGNORE INTO Dept (dept_id, dept_name) VALUES 
(1, '志愿服务部'),
(2, '社会福利部'),
(3, '环保宣传部'),
(4, '教育培训部');

-- 插入活动分类数据
INSERT IGNORE INTO ActivityCategory (category_id, category_name) VALUES 
(1, '社区服务'),
(2, '环境保护'),
(3, '教育培训'),
(4, '老人护理'),
(5, '儿童关爱');

-- 志愿者活动管理系统 - 数据库表结构调整脚本
-- 执行时间: 2025-12-22

-- ============================================================
-- 1. 为Activity表添加status字段（活动状态）
-- ============================================================
ALTER TABLE Activity ADD COLUMN status VARCHAR(20) DEFAULT 'active' COMMENT '活动状态: active(活跃), expired(已过期), closed(已关闭)';

-- ============================================================
-- 2. 为Activity表添加创建时间和修改时间字段（可选，便于审计）
-- ============================================================
ALTER TABLE Activity ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';
ALTER TABLE Activity ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间';

-- ============================================================
-- 3. 为User表添加创建时间字段（可选）
-- ============================================================
ALTER TABLE User ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间';

-- ============================================================
-- 4. 为Application表添加创建时间和修改时间字段（可选）
-- ============================================================
ALTER TABLE Application ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';
ALTER TABLE Application ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间';

-- ============================================================
-- 5. 为ApplicationStatusLog表添加创建时间字段（可选）
-- ============================================================
ALTER TABLE ApplicationStatusLog ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

-- ============================================================
-- 6. 创建索引以提高查询性能
-- ============================================================

-- Activity表的索引
CREATE INDEX idx_activity_status ON Activity(status);
CREATE INDEX idx_activity_time ON Activity(activity_time);
CREATE INDEX idx_activity_creator ON Activity(creator_id);

-- Application表的索引
CREATE INDEX idx_application_user ON Application(user_id);
CREATE INDEX idx_application_activity ON Application(activity_id);
CREATE INDEX idx_application_status ON Application(current_status);

-- ApplicationStatusLog表的索引
CREATE INDEX idx_status_log_application ON ApplicationStatusLog(application_id);

-- ============================================================
-- 7. 验证修改（查看Activity表的新字段）
-- ============================================================
-- DESC Activity;
-- SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT 
-- FROM INFORMATION_SCHEMA.COLUMNS 
-- WHERE TABLE_NAME='Activity' AND TABLE_SCHEMA=DATABASE();

-- ============================================
-- 腺样体面容智能筛查系统 - 数据库初始化脚本
-- ============================================

CREATE DATABASE IF NOT EXISTS face DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE face;

-- ============================================
-- 1. 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（bcrypt加密）',
    `name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '姓名',
    `avatar` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '头像URL',
    `email` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '邮箱',
    `phone` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号',
    `gender` TINYINT NOT NULL DEFAULT 0 COMMENT '性别：0-未知 1-男 2-女',
    `birthday` DATE DEFAULT NULL COMMENT '生日',
    `address` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '地址',
    `role` VARCHAR(20) NOT NULL DEFAULT 'doctor' COMMENT '角色：doctor-医生 admin-管理员',
    `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 2. 诊断任务表（异步诊断核心）
-- ============================================
CREATE TABLE IF NOT EXISTS `diagnosis_task` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务ID',
    `task_no` VARCHAR(64) NOT NULL COMMENT '任务编号（唯一幂等键）',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `image_path` VARCHAR(500) NOT NULL COMMENT '图片存储路径',
    `image_url` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '图片访问URL',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '任务状态：0-待处理 1-处理中 2-已完成 3-失败',
    `age` INT DEFAULT NULL COMMENT '患者年龄',
    `gender` TINYINT DEFAULT NULL COMMENT '患者性别：0-未知 1-男 2-女',
    `retry_count` INT NOT NULL DEFAULT 0 COMMENT '重试次数',
    `max_retry` INT NOT NULL DEFAULT 5 COMMENT '最大重试次数',
    `error_message` TEXT COMMENT '错误信息',
    `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `complete_time` DATETIME DEFAULT NULL COMMENT '完成时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_task_no` (`task_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='诊断任务表';

-- ============================================
-- 3. 诊断结果表
-- ============================================
CREATE TABLE IF NOT EXISTS `diagnosis_result` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '结果ID',
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT '关联任务ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `image_path` VARCHAR(500) NOT NULL COMMENT '图片路径',
    `image_url` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '图片URL',
    `is_gland_face` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否腺样体面容：0-否 1-是',
    `level` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '严重程度：轻度/中度/重度/非腺样体面容',
    `visualization_description` TEXT COMMENT '医学描述',
    `feature_vector_id` VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'Milvus向量ID',
    `reference_cases` TEXT COMMENT '引用的相似病例JSON',
    `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
    `test_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '检测时间',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_task_id` (`task_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_test_time` (`test_time`),
    KEY `idx_is_gland_face` (`is_gland_face`),
    KEY `idx_level` (`level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='诊断结果表';

-- ============================================
-- 4. 本地消息表（Outbox Pattern）
-- ============================================
CREATE TABLE IF NOT EXISTS `outbox_message` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    `aggregate_id` VARCHAR(64) NOT NULL COMMENT '聚合根ID（如任务ID）',
    `aggregate_type` VARCHAR(50) NOT NULL COMMENT '聚合类型',
    `message_type` VARCHAR(50) NOT NULL COMMENT '消息类型',
    `payload` TEXT NOT NULL COMMENT '消息内容JSON',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0-待发送 1-已发送 2-已确认 3-失败',
    `retry_count` INT NOT NULL DEFAULT 0 COMMENT '重试次数',
    `max_retry` INT NOT NULL DEFAULT 5 COMMENT '最大重试次数',
    `next_retry_time` DATETIME DEFAULT NULL COMMENT '下次重试时间',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_aggregate_id` (`aggregate_id`),
    KEY `idx_status` (`status`),
    KEY `idx_next_retry_time` (`next_retry_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='本地消息表（Outbox）';

-- ============================================
-- 5. 对话表
-- ============================================
CREATE TABLE IF NOT EXISTS `conversation` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '对话ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `title` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '对话标题',
    `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话表';

-- ============================================
-- 6. 消息表
-- ============================================
CREATE TABLE IF NOT EXISTS `message` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    `conversation_id` BIGINT UNSIGNED NOT NULL COMMENT '对话ID',
    `role` VARCHAR(20) NOT NULL COMMENT '角色：user/assistant',
    `content` TEXT NOT NULL COMMENT '消息内容',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_conversation_id` (`conversation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';

-- ============================================
-- 7. 诊断参数缓存表（三级缓存-MySQL层）
-- ============================================
CREATE TABLE IF NOT EXISTS `diagnosis_cache` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '缓存ID',
    `cache_key` VARCHAR(200) NOT NULL COMMENT '缓存键（年龄_性别组合等）',
    `cache_value` TEXT NOT NULL COMMENT '缓存值JSON',
    `expire_time` DATETIME DEFAULT NULL COMMENT '过期时间',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_cache_key` (`cache_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='诊断参数缓存表';

-- ============================================
-- 初始化默认管理员账户
-- 密码: admin123 (bcrypt: $2a$10$...)
-- ============================================
INSERT INTO `user` (`username`, `password`, `name`, `role`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '系统管理员', 'admin')
ON DUPLICATE KEY UPDATE `id`=`id`;

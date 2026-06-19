-- 修改日记表，添加用户ID字段
ALTER TABLE diary ADD COLUMN user_id BIGINT NOT NULL COMMENT '用户ID' AFTER id;

-- 添加外键约束
ALTER TABLE diary ADD CONSTRAINT fk_diary_user FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE;

-- 创建索引
CREATE INDEX idx_diary_user_id ON diary(user_id);
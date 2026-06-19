-- 创建数据库
CREATE DATABASE IF NOT EXISTS heart_shadow_diary CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE heart_shadow_diary;

-- 创建用户表
CREATE TABLE IF NOT EXISTS user (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(100) NOT NULL COMMENT '密码（加密存储）',
    avatar VARCHAR(255) COMMENT '头像URL',
    bio TEXT COMMENT '个人简介',
    followers INT DEFAULT 0 COMMENT '粉丝数',
    following INT DEFAULT 0 COMMENT '关注数',
    likes INT DEFAULT 0 COMMENT '获赞数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 创建心情表
CREATE TABLE IF NOT EXISTS mood (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '心情ID',
    name VARCHAR(50) NOT NULL COMMENT '心情名称',
    description TEXT COMMENT '心情描述',
    color VARCHAR(20) DEFAULT '#3498db' COMMENT '心情颜色',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='心情表';

-- 创建日记表
CREATE TABLE IF NOT EXISTS diary (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日记ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    title VARCHAR(100) NOT NULL COMMENT '日记标题',
    content TEXT COMMENT '日记内容',
    diary_date DATE NOT NULL COMMENT '日记日期',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日记表';

-- 创建日记照片表
CREATE TABLE IF NOT EXISTS diary_photo (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '照片ID',
    diary_id BIGINT NOT NULL COMMENT '关联的日记ID',
    photo_url VARCHAR(255) NOT NULL COMMENT '照片URL（阿里云OSS地址）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (diary_id) REFERENCES diary(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日记照片表';

-- 创建日记心情表
CREATE TABLE IF NOT EXISTS diary_mood (
    diary_id BIGINT NOT NULL COMMENT '日记ID',
    mood_id BIGINT NOT NULL COMMENT '心情ID',
    intensity INT DEFAULT 3 COMMENT '心情强度(1-5)',
    PRIMARY KEY (diary_id, mood_id),
    FOREIGN KEY (diary_id) REFERENCES diary(id) ON DELETE CASCADE,
    FOREIGN KEY (mood_id) REFERENCES mood(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日记心情表';

-- 创建AI会话表
CREATE TABLE IF NOT EXISTS ai_session (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '会话ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    session_name VARCHAR(100) COMMENT '会话名称',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI会话表';

-- 创建AI消息表
CREATE TABLE IF NOT EXISTS ai_message (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID',
    session_id BIGINT NOT NULL COMMENT '会话ID',
    sender VARCHAR(20) NOT NULL COMMENT '发送者（user/ai）',
    content TEXT NOT NULL COMMENT '消息内容',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (session_id) REFERENCES ai_session(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI消息表';

-- 创建绘梦表
CREATE TABLE IF NOT EXISTS dream_image (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '绘梦ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    prompt TEXT NOT NULL COMMENT '绘梦提示词',
    image_url VARCHAR(255) NOT NULL COMMENT '生成图片URL（阿里云OSS地址）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='绘梦表';

-- 创建人脸分析记录表
CREATE TABLE IF NOT EXISTS analysis_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    image_url VARCHAR(255) NOT NULL COMMENT '图片URL',
    analysis_date DATETIME NOT NULL COMMENT '分析日期时间',
    happy_count INT NOT NULL DEFAULT 0 COMMENT '开心表情数量',
    sad_count INT NOT NULL DEFAULT 0 COMMENT '悲伤表情数量',
    anger_count INT NOT NULL DEFAULT 0 COMMENT '愤怒表情数量',
    surprise_count INT NOT NULL DEFAULT 0 COMMENT '惊讶表情数量',
    fear_count INT NOT NULL DEFAULT 0 COMMENT '恐惧表情数量',
    neutral_count INT NOT NULL DEFAULT 0 COMMENT '中性表情数量',
    disgust_count INT NOT NULL DEFAULT 0 COMMENT '厌恶表情数量',
    dominant_expression VARCHAR(50) NOT NULL COMMENT '主要表情',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='人脸分析记录表';

-- 创建索引
CREATE INDEX idx_user_username ON user(username);
CREATE INDEX idx_diary_user_id ON diary(user_id);
CREATE INDEX idx_diary_diary_date ON diary(diary_date);
CREATE INDEX idx_diary_photo_diary_id ON diary_photo(diary_id);
CREATE INDEX idx_ai_session_user_id ON ai_session(user_id);
CREATE INDEX idx_ai_message_session_id ON ai_message(session_id);
CREATE INDEX idx_dream_image_user_id ON dream_image(user_id);
CREATE INDEX idx_analysis_record_user_id ON analysis_record(user_id);
CREATE INDEX idx_analysis_record_analysis_date ON analysis_record(analysis_date);
CREATE INDEX idx_analysis_record_dominant_expression ON analysis_record(dominant_expression);

-- 插入初始心情数据
INSERT INTO mood (name, description, color) VALUES 
('开心', '心情愉快，充满喜悦', '#f1c40f'),
('悲伤', '心情低落，感到难过', '#3498db'),
('愤怒', '情绪激动，感到生气', '#e74c3c'),
('惊讶', '感到意外，震惊', '#9b59b6'),
('自然', '心情平静，无特别情绪', '#2ecc71');

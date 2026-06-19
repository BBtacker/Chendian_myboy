-- 创建心情表
CREATE TABLE IF NOT EXISTS mood (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '心情ID',
    name VARCHAR(50) NOT NULL COMMENT '心情名称',
    description TEXT COMMENT '心情描述',
    color VARCHAR(20) DEFAULT '#3498db' COMMENT '心情颜色',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='心情表';

-- 插入初始心情数据
INSERT INTO mood (name, description, color) VALUES 
('开心', '心情愉快，充满喜悦', '#f1c40f'),
('悲伤', '心情低落，感到难过', '#3498db'),
('愤怒', '情绪激动，感到生气', '#e74c3c'),
('惊讶', '感到意外，震惊', '#9b59b6'),
('自然', '心情平静，无特别情绪', '#2ecc71');

-- 创建日记表
CREATE TABLE IF NOT EXISTS diary (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日记ID',
    title VARCHAR(100) NOT NULL COMMENT '日记标题',
    content TEXT COMMENT '日记内容',
    diary_date DATE NOT NULL COMMENT '日记日期',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日记表';

-- 创建日记照片表
CREATE TABLE IF NOT EXISTS diary_photo (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '照片ID',
    diary_id BIGINT NOT NULL COMMENT '关联的日记ID',
    photo_url VARCHAR(255) NOT NULL COMMENT '照片URL',
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

-- 创建索引
CREATE INDEX idx_diary_date ON diary(diary_date);
CREATE INDEX idx_diary_created_at ON diary(created_at);
CREATE INDEX idx_diary_photo_diary_id ON diary_photo(diary_id);

-- 保留原有人脸分析记录表，用于兼容旧数据
CREATE TABLE IF NOT EXISTS analysis_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
    image_url VARCHAR(255) NOT NULL COMMENT '图片URL',
    analysis_date DATETIME NOT NULL COMMENT '分析日期时间',
    happy_count INT NOT NULL DEFAULT 0 COMMENT '开心表情数量',
    sad_count INT NOT NULL DEFAULT 0 COMMENT '悲伤表情数量',
    anger_count INT NOT NULL DEFAULT 0 COMMENT '愤怒表情数量',
    surprise_count INT NOT NULL DEFAULT 0 COMMENT '惊讶表情数量',
    fear_count INT NOT NULL DEFAULT 0 COMMENT '恐惧表情数量',
    neutral_count INT NOT NULL DEFAULT 0 COMMENT '中性表情数量',
    disgust_count INT NOT NULL DEFAULT 0 COMMENT '厌恶表情数量',
    dominant_expression VARCHAR(50) NOT NULL COMMENT '主要表情'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='人脸分析记录表';

-- 创建索引
CREATE INDEX idx_analysis_date ON analysis_record(analysis_date);
CREATE INDEX idx_dominant_expression ON analysis_record(dominant_expression);

package com.lhy.face_recognition.entity;

import lombok.Data;
import java.time.LocalDateTime;

/**
 * AI会话实体类
 */
@Data
public class AiSession {
    private Long id;
    private Long userId;
    private String sessionName;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}
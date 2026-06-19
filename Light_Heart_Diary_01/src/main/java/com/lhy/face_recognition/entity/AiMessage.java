package com.lhy.face_recognition.entity;

import lombok.Data;
import java.time.LocalDateTime;

/**
 * AI消息实体类
 */
@Data
public class AiMessage {
    private Long id;
    private Long sessionId;
    private String sender;
    private String content;
    private LocalDateTime createdAt;
}
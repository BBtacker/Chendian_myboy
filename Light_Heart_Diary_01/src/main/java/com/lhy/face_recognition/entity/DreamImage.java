package com.lhy.face_recognition.entity;

import lombok.Data;
import java.time.LocalDateTime;

/**
 * 绘梦记录实体类
 */
@Data
public class DreamImage {
    private Long id;
    private Long userId;
    private String prompt;
    private String mood;
    private String imageUrl;
    private String originalImage;
    private LocalDateTime createdAt;
}
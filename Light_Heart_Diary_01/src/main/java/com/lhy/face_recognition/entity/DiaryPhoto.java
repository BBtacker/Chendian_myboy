package com.lhy.face_recognition.entity;

import lombok.Data;

import java.time.LocalDateTime;

/**
 * 日记照片实体类
 */
@Data
public class DiaryPhoto {
    /**
     * 照片ID
     */
    private Long id;
    
    /**
     * 关联的日记ID
     */
    private Long diaryId;
    
    /**
     * 照片URL
     */
    private String photoUrl;
    
    /**
     * 创建时间
     */
    private LocalDateTime createdAt;
}
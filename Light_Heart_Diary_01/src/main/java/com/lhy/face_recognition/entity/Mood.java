package com.lhy.face_recognition.entity;

import lombok.Data;

import java.time.LocalDateTime;

/**
 * 心情实体类
 */
@Data
public class Mood {
    /**
     * 心情ID
     */
    private Long id;
    
    /**
     * 心情名称
     */
    private String name;
    
    /**
     * 心情描述
     */
    private String description;
    
    /**
     * 心情颜色
     */
    private String color;
    
    /**
     * 创建时间
     */
    private LocalDateTime createdAt;
}
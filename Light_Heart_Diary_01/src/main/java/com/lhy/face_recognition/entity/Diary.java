package com.lhy.face_recognition.entity;

import lombok.Data;

import java.time.LocalDate;
import java.time.LocalDateTime;

/**
 * 日记实体类
 */
@Data
public class Diary {
    /**
     * 日记ID
     */
    private Long id;
    
    /**
     * 用户ID
     */
    private Long userId;
    
    /**
     * 日记标题
     */
    private String title;
    
    /**
     * 日记内容
     */
    private String content;
    
    /**
     * 日记日期
     */
    private LocalDate diaryDate;
    
    /**
     * 创建时间
     */
    private LocalDateTime createdAt;
    
    /**
     * 更新时间
     */
    private LocalDateTime updatedAt;
    
    /**
     * 心情
     */
    private String selectedMood;
}
package com.lhy.face_recognition.entity;

import lombok.Data;

/**
 * 日记心情关联实体类
 */
@Data
public class DiaryMood {
    /**
     * 日记ID
     */
    private Long diaryId;
    
    /**
     * 心情ID
     */
    private Long moodId;
    
    /**
     * 心情强度(1-5)
     */
    private Integer intensity;
}
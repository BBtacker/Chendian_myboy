package com.lhy.face_recognition.entity;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

import java.time.LocalDateTime;

/**
 * 人脸分析记录实体类
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
public class AnalysisRecord {
    /**
     * 记录ID
     */
    private Long id;
    
    /**
     * 图片URL
     */
    private String imageUrl;
    
    /**
     * 分析日期时间
     */
    private LocalDateTime analysisDate;
    
    /**
     * 开心表情数量
     */
    private Integer happyCount;
    
    /**
     * 悲伤表情数量
     */
    private Integer sadCount;
    
    /**
     * 愤怒表情数量
     */
    private Integer angerCount;
    
    /**
     * 惊讶表情数量
     */
    private Integer surpriseCount;
    
    /**
     * 恐惧表情数量
     */
    private Integer fearCount;
    
    /**
     * 中性表情数量
     */
    private Integer neutralCount;
    
    /**
     * 厌恶表情数量
     */
    private Integer disgustCount;
    
    /**
     * 蔑视表情数量
     */
    private Integer contemptCount;
    
    /**
     * 困倦表情数量
     */
    private Integer sleepyCount;
    
    /**
     * 主要表情
     */
    private String dominantExpression;
}

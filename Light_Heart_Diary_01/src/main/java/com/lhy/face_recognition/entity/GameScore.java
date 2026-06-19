package com.lhy.face_recognition.entity;

import lombok.Data;
import java.time.LocalDate;
import java.time.LocalDateTime;

/**
 * 游戏得分记录实体类
 */
@Data
public class GameScore {
    private Long id;
    private Long userId;
    private Integer score;
    private Integer highScore;
    private Integer positiveCount;
    private Integer negativeCount;
    private Integer combo;
    private Integer energy;
    private LocalDate gameDate;
    private LocalDateTime createdAt;
}
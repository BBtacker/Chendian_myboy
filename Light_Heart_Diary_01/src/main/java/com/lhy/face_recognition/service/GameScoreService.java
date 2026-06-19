package com.lhy.face_recognition.service;

import com.lhy.face_recognition.entity.GameScore;
import java.util.List;

public interface GameScoreService {
    GameScore saveGameScore(GameScore gameScore);
    List<GameScore> getGameScoresByUserId(Long userId);
    GameScore getBestScoreByUserId(Long userId);
    Integer getHighScoreByUserId(Long userId);
}
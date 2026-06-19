package com.lhy.face_recognition.controller;

import com.lhy.face_recognition.entity.GameScore;
import com.lhy.face_recognition.service.GameScoreService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/game")
@CrossOrigin(origins = {"http://localhost:5173", "http://localhost:5174"}, allowCredentials = "true")
public class GameScoreController {

    @Autowired
    private GameScoreService gameScoreService;

    /**
     * 保存游戏得分
     */
    @PostMapping("/score")
    public ResponseEntity<Map<String, Object>> saveScore(@RequestBody GameScore gameScore) {
        try {
            GameScore saved = gameScoreService.saveGameScore(gameScore);
            return ResponseEntity.ok(Map.of("success", true, "data", saved));
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(Map.of("success", false, "error", e.getMessage()));
        }
    }

    /**
     * 获取用户的所有游戏记录
     */
    @GetMapping("/scores/{userId}")
    public ResponseEntity<List<GameScore>> getScores(@PathVariable Long userId) {
        return ResponseEntity.ok(gameScoreService.getGameScoresByUserId(userId));
    }

    /**
     * 获取用户的最佳成绩
     */
    @GetMapping("/best/{userId}")
    public ResponseEntity<GameScore> getBestScore(@PathVariable Long userId) {
        GameScore best = gameScoreService.getBestScoreByUserId(userId);
        return best != null ? ResponseEntity.ok(best) : ResponseEntity.noContent().build();
    }

    /**
     * 获取用户的最高分
     */
    @GetMapping("/high-score/{userId}")
    public ResponseEntity<Map<String, Object>> getHighScore(@PathVariable Long userId) {
        Integer highScore = gameScoreService.getHighScoreByUserId(userId);
        return ResponseEntity.ok(Map.of("success", true, "highScore", highScore));
    }
}
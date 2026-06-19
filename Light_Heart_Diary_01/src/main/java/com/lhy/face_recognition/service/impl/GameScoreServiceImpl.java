package com.lhy.face_recognition.service.impl;

import com.lhy.face_recognition.entity.GameScore;
import com.lhy.face_recognition.mapper.GameScoreMapper;
import com.lhy.face_recognition.service.GameScoreService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class GameScoreServiceImpl implements GameScoreService {

    @Autowired
    private GameScoreMapper gameScoreMapper;

    @Override
    public GameScore saveGameScore(GameScore gameScore) {
        gameScoreMapper.insertGameScore(gameScore);
        return gameScore;
    }

    @Override
    public List<GameScore> getGameScoresByUserId(Long userId) {
        return gameScoreMapper.getGameScoresByUserId(userId);
    }

    @Override
    public GameScore getBestScoreByUserId(Long userId) {
        return gameScoreMapper.getBestScoreByUserId(userId);
    }

    @Override
    public Integer getHighScoreByUserId(Long userId) {
        return gameScoreMapper.getHighScoreByUserId(userId);
    }
}
package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.GameScore;
import org.apache.ibatis.annotations.*;

import java.util.List;

/**
 * 游戏得分Mapper接口
 */
@Mapper
public interface GameScoreMapper {

    @Insert("INSERT INTO game_score (user_id, score, high_score, positive_count, negative_count, combo, energy, game_date) " +
            "VALUES (#{userId}, #{score}, #{highScore}, #{positiveCount}, #{negativeCount}, #{combo}, #{energy}, #{gameDate})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertGameScore(GameScore gameScore);

    @Select("SELECT * FROM game_score WHERE user_id = #{userId} ORDER BY created_at DESC")
    List<GameScore> getGameScoresByUserId(Long userId);

    @Select("SELECT * FROM game_score WHERE user_id = #{userId} ORDER BY score DESC LIMIT 1")
    GameScore getBestScoreByUserId(Long userId);

    @Select("SELECT COALESCE(MAX(high_score), 0) FROM game_score WHERE user_id = #{userId}")
    Integer getHighScoreByUserId(Long userId);
}
package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.DiaryMood;
import org.apache.ibatis.annotations.*;

import java.util.List;

/**
 * 日记心情关联Mapper接口
 */
@Mapper
public interface DiaryMoodMapper {
    /**
     * 插入日记心情关联
     * @param diaryMood 日记心情关联实体
     * @return 影响的行数
     */
    @Insert("INSERT INTO diary_mood (diary_id, mood_id, intensity) VALUES (#{diaryId}, #{moodId}, #{intensity})")
    int insertDiaryMood(DiaryMood diaryMood);
    
    /**
     * 根据日记ID获取心情关联列表
     * @param diaryId 日记ID
     * @return 心情关联列表
     */
    @Select("SELECT * FROM diary_mood WHERE diary_id = #{diaryId}")
    List<DiaryMood> getMoodsByDiaryId(Long diaryId);
    
    /**
     * 根据日记ID删除所有心情关联
     * @param diaryId 日记ID
     * @return 影响的行数
     */
    @Delete("DELETE FROM diary_mood WHERE diary_id = #{diaryId}")
    int deleteMoodsByDiaryId(Long diaryId);
}
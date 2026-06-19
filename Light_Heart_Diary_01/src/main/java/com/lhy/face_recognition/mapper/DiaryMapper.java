package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.Diary;
import org.apache.ibatis.annotations.*;

import java.time.LocalDate;
import java.util.List;

/**
 * 日记Mapper接口
 */
@Mapper
public interface DiaryMapper {
    /**
     * 插入日记
     * @param diary 日记实体
     * @return 影响的行数
     */
    @Insert("INSERT INTO diary (user_id, title, content, diary_date, selected_mood, created_at, updated_at) VALUES (#{userId}, #{title}, #{content}, #{diaryDate}, #{selectedMood}, #{createdAt}, #{updatedAt})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertDiary(Diary diary);
    
    /**
     * 根据ID获取日记
     * @param id 日记ID
     * @return 日记实体
     */
    @Select("SELECT * FROM diary WHERE id = #{id}")
    Diary getDiaryById(Long id);
    
    /**
     * 根据用户ID获取所有日记
     * @param userId 用户ID
     * @return 日记列表
     */
    @Select("SELECT * FROM diary WHERE user_id = #{userId} ORDER BY diary_date DESC, created_at DESC")
    List<Diary> getDiariesByUserId(Long userId);
    
    /**
     * 根据用户ID和日期范围获取日记
     * @param userId 用户ID
     * @param startDate 开始日期
     * @param endDate 结束日期
     * @return 日记列表
     */
    @Select("SELECT * FROM diary WHERE user_id = #{userId} AND diary_date BETWEEN #{startDate} AND #{endDate} ORDER BY diary_date DESC, created_at DESC")
    List<Diary> getDiariesByUserIdAndDateRange(@Param("userId") Long userId, @Param("startDate") LocalDate startDate, @Param("endDate") LocalDate endDate);
    
    /**
     * 更新日记
     * @param diary 日记实体
     * @return 影响的行数
     */
    @Update("UPDATE diary SET title = #{title}, content = #{content}, diary_date = #{diaryDate}, selected_mood = #{selectedMood}, updated_at = #{updatedAt} WHERE id = #{id} AND user_id = #{userId}")
    int updateDiary(Diary diary);
    
    /**
     * 删除日记
     * @param id 日记ID
     * @param userId 用户ID
     * @return 影响的行数
     */
    @Delete("DELETE FROM diary WHERE id = #{id} AND user_id = #{userId}")
    int deleteDiary(@Param("id") Long id, @Param("userId") Long userId);
}
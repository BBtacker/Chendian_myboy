package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.Mood;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Select;

import java.util.List;

/**
 * 心情Mapper接口
 */
@Mapper
public interface MoodMapper {
    /**
     * 获取所有心情
     * @return 心情列表
     */
    @Select("SELECT * FROM mood ORDER BY id")
    List<Mood> getAllMoods();
    
    /**
     * 根据ID获取心情
     * @param id 心情ID
     * @return 心情
     */
    @Select("SELECT * FROM mood WHERE id = #{id}")
    Mood getMoodById(Long id);
    
    /**
     * 根据名称获取心情
     * @param name 心情名称
     * @return 心情
     */
    @Select("SELECT * FROM mood WHERE name = #{name}")
    Mood getMoodByName(String name);
}
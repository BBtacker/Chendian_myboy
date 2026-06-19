package com.lhy.face_recognition.service;

import com.lhy.face_recognition.entity.Mood;

import java.util.List;

/**
 * 心情Service接口
 */
public interface MoodService {
    /**
     * 获取所有心情
     * @return 心情列表
     */
    List<Mood> getAllMoods();
    
    /**
     * 根据ID获取心情
     * @param id 心情ID
     * @return 心情
     */
    Mood getMoodById(Long id);
    
    /**
     * 根据名称获取心情
     * @param name 心情名称
     * @return 心情
     */
    Mood getMoodByName(String name);
}
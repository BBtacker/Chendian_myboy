package com.lhy.face_recognition.service.impl;

import com.lhy.face_recognition.entity.Mood;
import com.lhy.face_recognition.mapper.MoodMapper;
import com.lhy.face_recognition.service.MoodService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * 心情Service实现类
 */
@Service
public class MoodServiceImpl implements MoodService {
    
    @Autowired
    private MoodMapper moodMapper;
    
    @Override
    public List<Mood> getAllMoods() {
        return moodMapper.getAllMoods();
    }
    
    @Override
    public Mood getMoodById(Long id) {
        return moodMapper.getMoodById(id);
    }
    
    @Override
    public Mood getMoodByName(String name) {
        return moodMapper.getMoodByName(name);
    }
}
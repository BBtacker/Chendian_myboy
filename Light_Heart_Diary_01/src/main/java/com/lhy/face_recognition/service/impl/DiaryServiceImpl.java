package com.lhy.face_recognition.service.impl;

import com.lhy.face_recognition.entity.Diary;
import com.lhy.face_recognition.entity.DiaryPhoto;
import com.lhy.face_recognition.entity.DiaryMood;
import com.lhy.face_recognition.mapper.DiaryMapper;
import com.lhy.face_recognition.mapper.DiaryPhotoMapper;
import com.lhy.face_recognition.mapper.DiaryMoodMapper;
import com.lhy.face_recognition.service.DiaryService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;

/**
 * 日记Service实现类
 */
@Service
public class DiaryServiceImpl implements DiaryService {
    
    @Autowired
    private DiaryMapper diaryMapper;
    
    @Autowired
    private DiaryPhotoMapper diaryPhotoMapper;
    
    @Autowired
    private DiaryMoodMapper diaryMoodMapper;
    
    @Override
    public Diary createDiary(Diary diary) {
        // 设置创建时间和更新时间
        LocalDateTime now = LocalDateTime.now();
        diary.setCreatedAt(now);
        diary.setUpdatedAt(now);
        
        // 插入日记
        diaryMapper.insertDiary(diary);
        
        return diary;
    }
    
    @Override
    public Diary getDiaryById(Long id, Long userId) {
        Diary diary = diaryMapper.getDiaryById(id);
        // 验证日记是否属于该用户
        if (diary != null && diary.getUserId().equals(userId)) {
            return diary;
        }
        return null;
    }
    
    @Override
    public List<Diary> getDiariesByUserId(Long userId) {
        return diaryMapper.getDiariesByUserId(userId);
    }
    
    @Override
    public List<Diary> getDiariesByUserIdAndDateRange(Long userId, LocalDate startDate, LocalDate endDate) {
        return diaryMapper.getDiariesByUserIdAndDateRange(userId, startDate, endDate);
    }
    
    @Override
    public Diary updateDiary(Diary diary) {
        // 设置更新时间
        diary.setUpdatedAt(LocalDateTime.now());
        
        // 更新日记
        int result = diaryMapper.updateDiary(diary);
        
        if (result > 0) {
            return diary;
        }
        return null;
    }
    
    @Override
    public boolean deleteDiary(Long id, Long userId) {
        // 验证日记是否属于该用户
        Diary diary = diaryMapper.getDiaryById(id);
        if (diary == null || !diary.getUserId().equals(userId)) {
            return false;
        }
        
        // 删除关联的照片
        diaryPhotoMapper.deletePhotosByDiaryId(id);
        
        // 删除关联的心情
        diaryMoodMapper.deleteMoodsByDiaryId(id);
        
        // 删除日记
        int result = diaryMapper.deleteDiary(id, userId);
        
        return result > 0;
    }
    
    @Override
    public DiaryPhoto addDiaryPhoto(DiaryPhoto diaryPhoto) {
        // 设置创建时间
        diaryPhoto.setCreatedAt(LocalDateTime.now());
        
        // 插入照片
        diaryPhotoMapper.insertDiaryPhoto(diaryPhoto);
        
        return diaryPhoto;
    }
    
    @Override
    public List<DiaryPhoto> getPhotosByDiaryId(Long diaryId) {
        return diaryPhotoMapper.getPhotosByDiaryId(diaryId);
    }
    
    @Override
    public boolean deletePhoto(Long photoId) {
        int result = diaryPhotoMapper.deletePhotoById(photoId);
        return result > 0;
    }
    
    @Override
    public boolean addDiaryMood(DiaryMood diaryMood) {
        int result = diaryMoodMapper.insertDiaryMood(diaryMood);
        return result > 0;
    }
    
    @Override
    public List<DiaryMood> getMoodsByDiaryId(Long diaryId) {
        return diaryMoodMapper.getMoodsByDiaryId(diaryId);
    }
    
    @Override
    public boolean deleteMoodsByDiaryId(Long diaryId) {
        int result = diaryMoodMapper.deleteMoodsByDiaryId(diaryId);
        return result > 0;
    }
}
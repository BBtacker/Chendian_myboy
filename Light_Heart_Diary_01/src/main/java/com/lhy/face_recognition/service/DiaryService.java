package com.lhy.face_recognition.service;

import com.lhy.face_recognition.entity.Diary;
import com.lhy.face_recognition.entity.DiaryPhoto;
import com.lhy.face_recognition.entity.DiaryMood;

import java.time.LocalDate;
import java.util.List;

/**
 * 日记Service接口
 */
public interface DiaryService {
    /**
     * 创建日记
     * @param diary 日记实体
     * @return 创建的日记
     */
    Diary createDiary(Diary diary);
    
    /**
     * 根据ID获取日记
     * @param id 日记ID
     * @param userId 用户ID
     * @return 日记实体
     */
    Diary getDiaryById(Long id, Long userId);
    
    /**
     * 根据用户ID获取所有日记
     * @param userId 用户ID
     * @return 日记列表
     */
    List<Diary> getDiariesByUserId(Long userId);
    
    /**
     * 根据用户ID和日期范围获取日记
     * @param userId 用户ID
     * @param startDate 开始日期
     * @param endDate 结束日期
     * @return 日记列表
     */
    List<Diary> getDiariesByUserIdAndDateRange(Long userId, LocalDate startDate, LocalDate endDate);
    
    /**
     * 更新日记
     * @param diary 日记实体
     * @return 更新后的日记
     */
    Diary updateDiary(Diary diary);
    
    /**
     * 删除日记
     * @param id 日记ID
     * @param userId 用户ID
     * @return 是否删除成功
     */
    boolean deleteDiary(Long id, Long userId);
    
    /**
     * 为日记添加照片
     * @param diaryPhoto 日记照片实体
     * @return 添加的照片
     */
    DiaryPhoto addDiaryPhoto(DiaryPhoto diaryPhoto);
    
    /**
     * 根据日记ID获取照片列表
     * @param diaryId 日记ID
     * @return 照片列表
     */
    List<DiaryPhoto> getPhotosByDiaryId(Long diaryId);
    
    /**
     * 删除照片
     * @param photoId 照片ID
     * @return 是否删除成功
     */
    boolean deletePhoto(Long photoId);
    
    /**
     * 为日记添加心情
     * @param diaryMood 日记心情关联实体
     * @return 是否添加成功
     */
    boolean addDiaryMood(DiaryMood diaryMood);
    
    /**
     * 根据日记ID获取心情列表
     * @param diaryId 日记ID
     * @return 心情关联列表
     */
    List<DiaryMood> getMoodsByDiaryId(Long diaryId);
    
    /**
     * 删除日记的所有心情关联
     * @param diaryId 日记ID
     * @return 是否删除成功
     */
    boolean deleteMoodsByDiaryId(Long diaryId);
}
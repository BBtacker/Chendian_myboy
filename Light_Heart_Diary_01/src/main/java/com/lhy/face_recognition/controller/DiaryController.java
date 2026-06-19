package com.lhy.face_recognition.controller;

import com.lhy.face_recognition.entity.Diary;
import com.lhy.face_recognition.entity.DiaryPhoto;
import com.lhy.face_recognition.entity.DiaryMood;
import com.lhy.face_recognition.entity.Mood;
import com.lhy.face_recognition.service.DiaryService;
import com.lhy.face_recognition.service.MoodService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDate;
import java.util.List;

/**
 * 日记Controller
 */
@RestController
@RequestMapping("/api/diary")
@CrossOrigin(origins = {"http://localhost:5173", "http://localhost:5174"}, allowCredentials = "true")
public class DiaryController {
    
    @Autowired
    private DiaryService diaryService;
    
    @Autowired
    private MoodService moodService;
    
    /**
     * 获取所有心情
     */
    @GetMapping("/moods")
    public ResponseEntity<List<Mood>> getAllMoods() {
        List<Mood> moods = moodService.getAllMoods();
        return ResponseEntity.ok(moods);
    }
    
    /**
     * 创建日记
     */
    @PostMapping
    public ResponseEntity<Diary> createDiary(@RequestBody Diary diary) {
        try {
            Diary createdDiary = diaryService.createDiary(diary);
            return ResponseEntity.status(HttpStatus.CREATED).body(createdDiary);
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(null);
        }
    }
    
    /**
     * 根据用户ID获取所有日记
     */
    @GetMapping("/user/{userId}")
    public ResponseEntity<List<Diary>> getDiariesByUserId(@PathVariable Long userId) {
        List<Diary> diaries = diaryService.getDiariesByUserId(userId);
        return ResponseEntity.ok(diaries);
    }
    
    /**
     * 根据ID获取日记
     */
    @GetMapping("/{id}/{userId}")
    public ResponseEntity<Diary> getDiaryById(@PathVariable Long id, @PathVariable Long userId) {
        Diary diary = diaryService.getDiaryById(id, userId);
        if (diary == null) {
            return ResponseEntity.notFound().build();
        }
        return ResponseEntity.ok(diary);
    }
    
    /**
     * 根据用户ID和日期范围获取日记
     */
    @GetMapping("/user/{userId}/date-range")
    public ResponseEntity<List<Diary>> getDiariesByUserIdAndDateRange(
            @PathVariable Long userId,
            @RequestParam @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate startDate,
            @RequestParam @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate endDate) {
        List<Diary> diaries = diaryService.getDiariesByUserIdAndDateRange(userId, startDate, endDate);
        return ResponseEntity.ok(diaries);
    }
    
    /**
     * 更新日记
     */
    @PutMapping("/{id}")
    public ResponseEntity<Diary> updateDiary(@PathVariable Long id, @RequestBody Diary diary) {
        try {
            // 确保ID一致
            diary.setId(id);
            Diary updatedDiary = diaryService.updateDiary(diary);
            if (updatedDiary == null) {
                return ResponseEntity.notFound().build();
            }
            return ResponseEntity.ok(updatedDiary);
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(null);
        }
    }
    
    /**
     * 删除日记
     */
    @DeleteMapping("/{id}/{userId}")
    public ResponseEntity<Boolean> deleteDiary(@PathVariable Long id, @PathVariable Long userId) {
        boolean result = diaryService.deleteDiary(id, userId);
        if (result) {
            return ResponseEntity.ok(true);
        } else {
            return ResponseEntity.notFound().build();
        }
    }
    
    /**
     * 为日记添加照片
     */
    @PostMapping("/{id}/photos")
    public ResponseEntity<DiaryPhoto> addDiaryPhoto(@PathVariable Long id, @RequestBody DiaryPhoto diaryPhoto) {
        try {
            // 确保日记ID一致
            diaryPhoto.setDiaryId(id);
            DiaryPhoto addedPhoto = diaryService.addDiaryPhoto(diaryPhoto);
            return ResponseEntity.status(HttpStatus.CREATED).body(addedPhoto);
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(null);
        }
    }
    
    /**
     * 根据日记ID获取照片列表
     */
    @GetMapping("/{id}/photos")
    public ResponseEntity<List<DiaryPhoto>> getPhotosByDiaryId(@PathVariable Long id) {
        List<DiaryPhoto> photos = diaryService.getPhotosByDiaryId(id);
        return ResponseEntity.ok(photos);
    }
    
    /**
     * 删除照片
     */
    @DeleteMapping("/photos/{photoId}")
    public ResponseEntity<Boolean> deletePhoto(@PathVariable Long photoId) {
        boolean result = diaryService.deletePhoto(photoId);
        if (result) {
            return ResponseEntity.ok(true);
        } else {
            return ResponseEntity.notFound().build();
        }
    }
    
    /**
     * 为日记添加心情
     */
    @PostMapping("/{id}/moods")
    public ResponseEntity<Boolean> addDiaryMood(@PathVariable Long id, @RequestBody DiaryMood diaryMood) {
        try {
            // 确保日记ID一致
            diaryMood.setDiaryId(id);
            boolean result = diaryService.addDiaryMood(diaryMood);
            return ResponseEntity.status(HttpStatus.CREATED).body(result);
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(false);
        }
    }
    
    /**
     * 根据日记ID获取心情列表
     */
    @GetMapping("/{id}/moods")
    public ResponseEntity<List<DiaryMood>> getMoodsByDiaryId(@PathVariable Long id) {
        List<DiaryMood> moods = diaryService.getMoodsByDiaryId(id);
        return ResponseEntity.ok(moods);
    }
    
    /**
     * 删除日记的所有心情关联
     */
    @DeleteMapping("/{id}/moods")
    public ResponseEntity<Boolean> deleteMoodsByDiaryId(@PathVariable Long id) {
        boolean result = diaryService.deleteMoodsByDiaryId(id);
        return ResponseEntity.ok(result);
    }
}
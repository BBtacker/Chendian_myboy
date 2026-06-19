package com.lhy.face_recognition.controller;

import com.lhy.face_recognition.entity.AnalysisRecord;
import com.lhy.face_recognition.service.AnalysisRecordService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.time.LocalDateTime;
import java.util.List;

/**
 * 人脸分析记录Controller
 */
@RestController
@RequestMapping("/api/analysis")
@CrossOrigin(origins = {"http://localhost:5173", "http://localhost:5174"}, allowCredentials = "true") // 允许前端跨域访问并支持凭证
public class AnalysisController {
    
    @Autowired
    private AnalysisRecordService analysisRecordService;
    
    /**
     * 获取所有分析记录
     */
    @GetMapping("/records")
    public ResponseEntity<List<AnalysisRecord>> getAllRecords() {
        List<AnalysisRecord> records = analysisRecordService.getAllRecords();
        return ResponseEntity.ok(records);
    }
    
    /**
     * 获取最新的分析记录
     */
    @GetMapping("/records/latest")
    public ResponseEntity<List<AnalysisRecord>> getLatestRecords(@RequestParam(defaultValue = "10") int limit) {
        List<AnalysisRecord> records = analysisRecordService.getLatestRecords(limit);
        return ResponseEntity.ok(records);
    }
    
    /**
     * 根据ID获取分析记录
     */
    @GetMapping("/records/{id}")
    public ResponseEntity<AnalysisRecord> getRecordById(@PathVariable Long id) {
        AnalysisRecord record = analysisRecordService.getRecordById(id);
        if (record == null) {
            return ResponseEntity.notFound().build();
        }
        return ResponseEntity.ok(record);
    }
    
    /**
     * 根据日期范围获取分析记录
     */
    @GetMapping("/records/date-range")
    public ResponseEntity<List<AnalysisRecord>> getRecordsByDateRange(
            @RequestParam @DateTimeFormat(pattern = "yyyy-MM-dd HH:mm:ss") LocalDateTime startDate,
            @RequestParam @DateTimeFormat(pattern = "yyyy-MM-dd HH:mm:ss") LocalDateTime endDate) {
        
        List<AnalysisRecord> records = analysisRecordService.getRecordsByDateRange(startDate, endDate);
        return ResponseEntity.ok(records);
    }
    
    /**
     * 根据主要表情获取分析记录
     */
    @GetMapping("/records/dominant-expression/{expression}")
    public ResponseEntity<List<AnalysisRecord>> getRecordsByDominantExpression(@PathVariable String expression) {
        List<AnalysisRecord> records = analysisRecordService.getRecordsByDominantExpression(expression);
        return ResponseEntity.ok(records);
    }
    
    /**
     * 上传图片进行分析
     */
    @PostMapping("/analyze")
    public ResponseEntity<Object> analyzeImage(@RequestParam("file") MultipartFile file) {
        try {
            // 检查文件是否为空
            if (file.isEmpty()) {
                return ResponseEntity.badRequest().body("未提供有效的图像文件");
            }
            
            // 获取文件名和文件内容
            String fileName = file.getOriginalFilename();
            byte[] imageBytes = file.getBytes();
            
            // 调用服务层进行分析
            AnalysisRecord record = analysisRecordService.analyzeImage(imageBytes, fileName);
            
            return ResponseEntity.ok(record);
        } catch (IOException e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("图片处理失败: " + e.getMessage());
        } catch (RuntimeException e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("分析失败: " + e.getMessage());
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("系统错误: " + e.getMessage());
        }
    }
    
    /**
     * 删除分析记录
     */
    @DeleteMapping("/records/{id}")
    public ResponseEntity<Boolean> deleteRecord(@PathVariable Long id) {
        boolean result = analysisRecordService.deleteRecord(id);
        if (result) {
            return ResponseEntity.ok(true);
        } else {
            return ResponseEntity.notFound().build();
        }
    }
    
    /**
     * 更新分析记录的表情数据（手动录入）
     */
    @PutMapping("/records/{id}/emotions")
    public ResponseEntity<AnalysisRecord> updateEmotions(@PathVariable Long id, @RequestBody EmotionUpdateRequest request) {
        try {
            // 获取现有记录
            AnalysisRecord record = analysisRecordService.getRecordById(id);
            if (record == null) {
                return ResponseEntity.notFound().build();
            }
            
            // 更新表情数据
            record.setHappyCount(request.getHappyCount());
            record.setSadCount(request.getSadCount());
            record.setAngerCount(request.getAngryCount());
            record.setSurpriseCount(request.getSurpriseCount());
            record.setNeutralCount(request.getNaturalCount());
            record.setDominantExpression(request.getDominantExpression());
            
            // 保存更新后的记录
            AnalysisRecord updatedRecord = analysisRecordService.saveRecord(record);
            
            return ResponseEntity.ok(updatedRecord);
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(null);
        }
    }
    
    /**
     * 表情更新请求体
     */
    public static class EmotionUpdateRequest {
        private int happyCount;
        private int sadCount;
        private int angryCount;
        private int surpriseCount;
        private int naturalCount;
        private String dominantExpression;
        
        // Getters and setters
        public int getHappyCount() {
            return happyCount;
        }
        
        public void setHappyCount(int happyCount) {
            this.happyCount = happyCount;
        }
        
        public int getSadCount() {
            return sadCount;
        }
        
        public void setSadCount(int sadCount) {
            this.sadCount = sadCount;
        }
        
        public int getAngryCount() {
            return angryCount;
        }
        
        public void setAngryCount(int angryCount) {
            this.angryCount = angryCount;
        }
        
        public int getSurpriseCount() {
            return surpriseCount;
        }
        
        public void setSurpriseCount(int surpriseCount) {
            this.surpriseCount = surpriseCount;
        }
        
        public int getNaturalCount() {
            return naturalCount;
        }
        
        public void setNaturalCount(int naturalCount) {
            this.naturalCount = naturalCount;
        }
        
        public String getDominantExpression() {
            return dominantExpression;
        }
        
        public void setDominantExpression(String dominantExpression) {
            this.dominantExpression = dominantExpression;
        }
    }
}

package com.lhy.face_recognition.service.impl;

import com.lhy.face_recognition.entity.AnalysisRecord;
import com.lhy.face_recognition.mapper.AnalysisRecordMapper;
import com.lhy.face_recognition.service.AnalysisRecordService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

/**
 * 人脸分析记录Service实现类
 */
@Service
public class AnalysisRecordServiceImpl implements AnalysisRecordService {
    
    @Autowired
    private AnalysisRecordMapper analysisRecordMapper;
    
    @Override
    public List<AnalysisRecord> getAllRecords() {
        return analysisRecordMapper.getAllRecords();
    }
    
    @Override
    public List<AnalysisRecord> getLatestRecords(int limit) {
        return analysisRecordMapper.getLatestRecords(limit);
    }
    
    @Override
    public AnalysisRecord getRecordById(Long id) {
        return analysisRecordMapper.getRecordById(id);
    }
    
    @Override
    public List<AnalysisRecord> getRecordsByDateRange(LocalDateTime startDate, LocalDateTime endDate) {
        return analysisRecordMapper.getRecordsByDateRange(startDate, endDate);
    }
    
    @Override
    public List<AnalysisRecord> getRecordsByDominantExpression(String dominantExpression) {
        return analysisRecordMapper.getRecordsByDominantExpression(dominantExpression);
    }
    
    @Override
    public AnalysisRecord saveRecord(AnalysisRecord record) {
        // 如果是新记录，设置分析日期
        if (record.getId() == null) {
            record.setAnalysisDate(LocalDateTime.now());
        }
        
        // 插入记录
        analysisRecordMapper.insertRecord(record);
        
        return record;
    }
    
    @Override
    public boolean deleteRecord(Long id) {
        int result = analysisRecordMapper.deleteRecord(id);
        return result > 0;
    }
    
    @Override
    public AnalysisRecord analyzeImage(byte[] imageBytes, String fileName) {
        // 直接保存图片信息，不调用Python服务
        AnalysisRecord record = new AnalysisRecord();
        record.setAnalysisDate(LocalDateTime.now());
        record.setImageUrl("/images/" + fileName);
        
        // 初始化表情计数为0（前端会手动录入）
        record.setHappyCount(0);
        record.setSadCount(0);
        record.setAngerCount(0);
        record.setSurpriseCount(0);
        record.setFearCount(0);
        record.setNeutralCount(0);
        record.setDisgustCount(0);
        record.setContemptCount(0);
        record.setSleepyCount(0);
        record.setDominantExpression("natural");
        
        return saveRecord(record);
    }
}

package com.lhy.face_recognition.service;

import com.lhy.face_recognition.entity.AnalysisRecord;

import java.time.LocalDateTime;
import java.util.List;

/**
 * 人脸分析记录Service接口
 */
public interface AnalysisRecordService {
    /**
     * 获取所有分析记录
     * @return 分析记录列表
     */
    List<AnalysisRecord> getAllRecords();
    
    /**
     * 获取最新的分析记录
     * @param limit 限制数量
     * @return 分析记录列表
     */
    List<AnalysisRecord> getLatestRecords(int limit);
    
    /**
     * 根据ID获取分析记录
     * @param id 记录ID
     * @return 分析记录
     */
    AnalysisRecord getRecordById(Long id);
    
    /**
     * 根据日期范围获取分析记录
     * @param startDate 开始日期
     * @param endDate 结束日期
     * @return 分析记录列表
     */
    List<AnalysisRecord> getRecordsByDateRange(LocalDateTime startDate, LocalDateTime endDate);
    
    /**
     * 根据主要表情获取分析记录
     * @param dominantExpression 主要表情
     * @return 分析记录列表
     */
    List<AnalysisRecord> getRecordsByDominantExpression(String dominantExpression);
    
    /**
     * 保存分析记录
     * @param record 分析记录
     * @return 保存后的分析记录
     */
    AnalysisRecord saveRecord(AnalysisRecord record);
    
    /**
     * 删除分析记录
     * @param id 记录ID
     * @return 是否删除成功
     */
    boolean deleteRecord(Long id);
    
    /**
     * 上传图片并分析表情（预留接口）
     * @param imageBytes 图片字节数组
     * @param fileName 文件名
     * @return 分析结果记录
     */
    AnalysisRecord analyzeImage(byte[] imageBytes, String fileName);
}

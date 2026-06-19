package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.AnalysisRecord;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

import java.time.LocalDateTime;
import java.util.List;

/**
 * 人脸分析记录Mapper接口
 */
@Mapper
public interface AnalysisRecordMapper {
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
    List<AnalysisRecord> getLatestRecords(@Param("limit") int limit);
    
    /**
     * 根据ID获取分析记录
     * @param id 记录ID
     * @return 分析记录
     */
    AnalysisRecord getRecordById(@Param("id") Long id);
    
    /**
     * 根据日期范围获取分析记录
     * @param startDate 开始日期
     * @param endDate 结束日期
     * @return 分析记录列表
     */
    List<AnalysisRecord> getRecordsByDateRange(@Param("startDate") LocalDateTime startDate, 
                                              @Param("endDate") LocalDateTime endDate);
    
    /**
     * 根据主要表情获取分析记录
     * @param dominantExpression 主要表情
     * @return 分析记录列表
     */
    List<AnalysisRecord> getRecordsByDominantExpression(@Param("dominantExpression") String dominantExpression);
    
    /**
     * 插入分析记录
     * @param record 分析记录
     * @return 影响行数
     */
    int insertRecord(AnalysisRecord record);
    
    /**
     * 删除分析记录
     * @param id 记录ID
     * @return 影响行数
     */
    int deleteRecord(@Param("id") Long id);
}

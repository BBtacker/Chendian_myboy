package com.long67.facetest.service.Impl;

import com.long67.facetest.mapper.StatisticsMapper;
import com.long67.facetest.service.StatisticsService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class StatisticsServiceImpl implements StatisticsService {
    
    @Autowired
    private StatisticsMapper statisticsMapper;
    
    @Override
    public Map<String, Object> getStatisticsOverview(Integer userId, LocalDate startDate, LocalDate endDate) {
        Map<String, Object> result = new HashMap<>();
        
        // 鑾峰彇鎬绘娴嬫鏁?        Integer totalTests = statisticsMapper.getTotalTests(userId, startDate, endDate);
        result.put("totalTests", totalTests != null ? totalTests : 0);
        
        // 鑾峰彇鑵轰綋闈㈠娆℃暟
        Integer glandFaceCount = statisticsMapper.getGlandFaceCount(userId, startDate, endDate);
        result.put("glandFaceCount", glandFaceCount != null ? glandFaceCount : 0);
        
        // 鑾峰彇闈炶吅浣撻潰瀹规鏁?        Integer nonGlandFaceCount = totalTests != null && glandFaceCount != null ? 
                totalTests - glandFaceCount : 0;
        result.put("nonGlandFaceCount", nonGlandFaceCount);
        
        // 璁＄畻鍑嗙‘鐜?(鎵€鏈夋娴嬬疆淇″害鐨勫钩鍧囧€?
        Double averageConfidence = statisticsMapper.getAverageConfidence(userId, startDate, endDate);
        double accuracyRate = averageConfidence != null ? averageConfidence : 0;
        result.put("accuracyRate", accuracyRate);
        
        // 鑾峰彇瓒嬪娍鏁版嵁
        List<Map<String, Object>> trendData = statisticsMapper.getTrendData(userId, startDate, endDate);
        result.put("trendData", trendData);
        
        // 鑾峰彇绛夌骇鍒嗗竷鏁版嵁
        List<Map<String, Object>> levelData = statisticsMapper.getLevelData(userId, startDate, endDate);
        result.put("levelData", levelData);
        
        return result;
    }
    
    @Override
    public List<Map<String, Object>> getStatisticsDetail(Integer userId, LocalDate startDate, LocalDate endDate) {
        return statisticsMapper.getDetailData(userId, startDate, endDate);
    }
    
    @Override
    public Double getAverageConfidence(Integer userId, LocalDate startDate, LocalDate endDate) {
        return statisticsMapper.getAverageConfidence(userId, startDate, endDate);
    }
}
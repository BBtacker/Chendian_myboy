package com.long67.facetest.service;

import java.time.LocalDate;
import java.util.List;
import java.util.Map;

public interface StatisticsService {
    
    /**
     * 鑾峰彇缁熻姒傝鏁版嵁
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 缁熻姒傝鏁版嵁
     */
    Map<String, Object> getStatisticsOverview(Integer userId, LocalDate startDate, LocalDate endDate);
    
    /**
     * 鑾峰彇璇︾粏缁熻鏁版嵁
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 璇︾粏缁熻鏁版嵁
     */
    List<Map<String, Object>> getStatisticsDetail(Integer userId, LocalDate startDate, LocalDate endDate);

    /**
     * 鑾峰彇骞冲潎缃俊搴?     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 骞冲潎缃俊搴?     */
    Double getAverageConfidence(Integer userId, LocalDate startDate, LocalDate endDate);
}
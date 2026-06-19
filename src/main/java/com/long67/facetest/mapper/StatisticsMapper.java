package com.long67.facetest.mapper;

import org.apache.ibatis.annotations.MapKey;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

import java.time.LocalDate;
import java.util.List;
import java.util.Map;

@Mapper
public interface StatisticsMapper {
    
    /**
     * 鑾峰彇鎬绘娴嬫鏁?     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 鎬绘娴嬫鏁?     */
    Integer getTotalTests(@Param("userId") Integer userId, 
                         @Param("startDate") LocalDate startDate, 
                         @Param("endDate") LocalDate endDate);
    
    /**
     * 鑾峰彇鑵轰綋闈㈠娆℃暟
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 鑵轰綋闈㈠娆℃暟
     */
    Integer getGlandFaceCount(@Param("userId") Integer userId, 
                             @Param("startDate") LocalDate startDate, 
                             @Param("endDate") LocalDate endDate);
    
    /**
     * 鑾峰彇骞冲潎缃俊搴?     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 骞冲潎缃俊搴?     */
    Double getAverageConfidence(@Param("userId") Integer userId,
                               @Param("startDate") LocalDate startDate,
                               @Param("endDate") LocalDate endDate);
    
    /**
     * 鑾峰彇瓒嬪娍鏁版嵁
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 瓒嬪娍鏁版嵁
     */
    @MapKey("date")
    List<Map<String, Object>> getTrendData(@Param("userId") Integer userId, 
                                          @Param("startDate") LocalDate startDate, 
                                          @Param("endDate") LocalDate endDate);
    
    /**
     * 鑾峰彇绛夌骇鍒嗗竷鏁版嵁
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 绛夌骇鍒嗗竷鏁版嵁
     */
    @MapKey("level")
    List<Map<String, Object>> getLevelData(@Param("userId") Integer userId, 
                                          @Param("startDate") LocalDate startDate, 
                                          @Param("endDate") LocalDate endDate);
    
    /**
     * 鑾峰彇璇︾粏缁熻鏁版嵁
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 璇︾粏缁熻鏁版嵁
     */
    @MapKey("date")
    List<Map<String, Object>> getDetailData(@Param("userId") Integer userId, 
                                           @Param("startDate") LocalDate startDate, 
                                           @Param("endDate") LocalDate endDate);
}
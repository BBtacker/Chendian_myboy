package com.long67.facetest.mapper;

import com.long67.facetest.pojo.testResult;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

import java.time.LocalDateTime;
import java.util.List;

@Mapper
public interface TestResultMapper {
    
    /**
     * 鐢ㄦ埛鏌ヨ鑷繁鐨勬娴嬬粨鏋滐紙鏀寔鏉′欢妫€绱級
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @return 妫€娴嬬粨鏋滃垪琛?
     */
    List<testResult> getUserTestResults(@Param("userId") Integer userId,
                                        @Param("startDate") LocalDateTime startDate,
                                        @Param("endDate") LocalDateTime endDate,
                                        @Param("isGlandFace") Boolean isGlandFace,
                                        @Param("level") String level);
    
    /**
     * 鐢ㄦ埛鍒犻櫎鑷繁鐨勬娴嬬粨鏋?
     * @param resultId 妫€娴嬬粨鏋淚D
     * @param userId 鐢ㄦ埛ID
     * @return 鍒犻櫎璁板綍鏁?
     */
    int deleteUserTestResult(@Param("resultId") Integer resultId, @Param("userId") Integer userId);
    
    /**
     * 鐢ㄦ埛鎵归噺鍒犻櫎妫€娴嬬粨鏋?
     * @param ids 妫€娴嬬粨鏋淚D鍒楄〃
     * @param userId 鐢ㄦ埛ID
     * @return 鍒犻櫎璁板綍鏁?
     */
    int deleteBatchUserTestResults(@Param("ids") List<Integer> ids, @Param("userId") Integer userId);

    /**
     * 鏌ヨ鐢ㄦ埛妫€娴嬬粨鏋滅敤浜庡鍑?
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @return 妫€娴嬬粨鏋滃垪琛?
     */
    List<testResult> getTestResultsForExport(@Param("userId") Integer userId,
                                             @Param("startDate") LocalDateTime startDate,
                                             @Param("endDate") LocalDateTime endDate,
                                             @Param("isGlandFace") Boolean isGlandFace,
                                             @Param("level") String level);
    
    /**
     * 鐢ㄦ埛鏇存柊鑷繁鐨勬娴嬬粨鏋?
     * @param testResult 妫€娴嬬粨鏋滃璞?
     * @param userId 鐢ㄦ埛ID
     * @return 鏇存柊璁板綍鏁?
     */
    int updateUserTestResult(@Param("testResult") testResult testResult, @Param("userId") Integer userId);
}
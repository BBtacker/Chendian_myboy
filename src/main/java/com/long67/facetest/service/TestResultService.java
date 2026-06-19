package com.long67.facetest.service;

import com.github.pagehelper.PageInfo;
import com.long67.facetest.pojo.testResult;
import jakarta.servlet.http.HttpServletResponse;

import java.time.LocalDateTime;
import java.util.List;

public interface TestResultService {
    
    /**
     * 鐢ㄦ埛鏌ョ湅鑷繁鐨勬娴嬬粨鏋滐紙鏀寔鍒嗛〉鍜屾潯浠舵绱級
     * @param userId 鐢ㄦ埛ID
     * @param page 椤电爜
     * @param pageSize 姣忛〉澶у皬
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @return 妫€娴嬬粨鏋滃垎椤典俊鎭?
     */
    PageInfo<testResult> getUserTestResults(Integer userId, Integer page, Integer pageSize,
                                            LocalDateTime startDate, LocalDateTime endDate,
                                            Boolean isGlandFace, String level);
    
    /**
     * 鐢ㄦ埛鍒犻櫎鑷繁鐨勬娴嬬粨鏋?
     * @param resultId 妫€娴嬬粨鏋淚D
     * @param userId 鐢ㄦ埛ID
     * @return 鍒犻櫎鏄惁鎴愬姛
     */
    boolean deleteUserTestResult(Integer resultId, Integer userId);
    
    /**
     * 鐢ㄦ埛鎵归噺鍒犻櫎妫€娴嬬粨鏋?
     * @param ids 妫€娴嬬粨鏋淚D鍒楄〃
     * @param userId 鐢ㄦ埛ID
     * @return 鍒犻櫎鐨勮褰曟暟
     */
    int deleteBatchUserTestResults(List<Integer> ids, Integer userId);
    
    /**
     * 瀵煎嚭鐢ㄦ埛妫€娴嬬粨鏋滃埌Excel
     * @param response HTTP鍝嶅簲瀵硅薄
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @throws Exception 瀵煎嚭寮傚父
     */
    void exportTestResultsToExcel(HttpServletResponse response, Integer userId, LocalDateTime startDate, LocalDateTime endDate,
                                  Boolean isGlandFace, String level) throws Exception;
    
    /**
     * 瀵煎嚭鐢ㄦ埛妫€娴嬬粨鏋滃埌PDF
     * @param response HTTP鍝嶅簲瀵硅薄
     * @param userId 鐢ㄦ埛ID
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @throws Exception 瀵煎嚭寮傚父
     */
    void exportTestResultsToPDF(HttpServletResponse response, Integer userId, LocalDateTime startDate, LocalDateTime endDate,
                                Boolean isGlandFace, String level) throws Exception;
    
    /**
     * 鐢ㄦ埛鏇存柊鑷繁鐨勬娴嬬粨鏋?
     * @param testResult 妫€娴嬬粨鏋滃璞?
     * @param userId 鐢ㄦ埛ID
     * @return 鏇存柊鏄惁鎴愬姛
     */
    boolean updateUserTestResult(testResult testResult, Integer userId);
}
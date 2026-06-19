package com.long67.facetest.controller;

import com.github.pagehelper.PageInfo;
import com.long67.facetest.pojo.Result;
import com.long67.facetest.pojo.testResult;
import com.long67.facetest.service.TestResultService;
import com.long67.facetest.utils.UserThreadLocal;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.web.bind.annotation.*;

import jakarta.servlet.http.HttpServletResponse;
import java.time.LocalDateTime;
import java.util.List;

@RestController
@RequestMapping("/testResult")
public class TestResultController {
    
    @Autowired
    private TestResultService testResultService;
    
    /**
     * 鐢ㄦ埛鏌ョ湅鑷繁鐨勬娴嬬粨鏋滐紙鏀寔鍒嗛〉鍜屾潯浠舵绱級
     * @param page 椤电爜
     * @param pageSize 姣忛〉澶у皬
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     * @return 妫€娴嬬粨鏋滃垪琛?
     */
    @GetMapping("/result")
    public Result getUserTestResults(@RequestParam(defaultValue = "1") Integer page,
                                     @RequestParam(defaultValue = "10") Integer pageSize,
                                     @DateTimeFormat(pattern = "yyyy-MM-dd HH:mm:ss") @RequestParam(required = false) LocalDateTime startDate,
                                     @DateTimeFormat(pattern = "yyyy-MM-dd HH:mm:ss") @RequestParam(required = false) LocalDateTime endDate,
                                     @RequestParam(required = false) Boolean isGlandFace,
                                     @RequestParam(required = false) String level) {
        try {
            // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
            Integer userId = UserThreadLocal.getUserId();
            PageInfo<testResult> pageInfo = testResultService.getUserTestResults(userId, page, pageSize, 
                    startDate, endDate, isGlandFace, level);
            return Result.success(pageInfo);
        } catch (Exception e) {
            return Result.error("鏌ヨ澶辫触: " + e.getMessage());
        }
    }
    
    /**
     * 鐢ㄦ埛鍒犻櫎鑷繁鐨勬娴嬬粨鏋?
     * @param resultId 妫€娴嬬粨鏋淚D
     * @return 鍒犻櫎缁撴灉
     */
    @DeleteMapping("/{resultId}")
    public Result deleteUserTestResult(@PathVariable Integer resultId) {
        try {
            // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
            Integer userId = UserThreadLocal.getUserId();
            boolean success = testResultService.deleteUserTestResult(resultId, userId);
            if (success) {
                return Result.success("鍒犻櫎鎴愬姛");
            } else {
                return Result.error("鍒犻櫎澶辫触锛岃褰曚笉瀛樺湪鎴栨棤鏉冮檺");
            }
        } catch (Exception e) {
            return Result.error("鍒犻櫎澶辫触: " + e.getMessage());
        }
    }
    
    /**
     * 鐢ㄦ埛鎵归噺鍒犻櫎妫€娴嬬粨鏋?
     * @param ids 妫€娴嬬粨鏋淚D鍒楄〃
     * @return 鍒犻櫎缁撴灉
     */
    @DeleteMapping("/batch")
    public Result deleteBatchUserTestResults(@RequestBody List<Integer> ids) {
        try {
            // 妫€鏌ュ弬鏁?
            if (ids == null || ids.isEmpty()) {
                return Result.error("璇烽€夋嫨瑕佸垹闄ょ殑璁板綍");
            }
            
            // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
            Integer userId = UserThreadLocal.getUserId();
            int deletedCount = testResultService.deleteBatchUserTestResults(ids, userId);
            return Result.success("鎴愬姛鍒犻櫎 " + deletedCount + " 鏉¤褰?);
        } catch (Exception e) {
            return Result.error("鎵归噺鍒犻櫎澶辫触: " + e.getMessage());
        }
    }
    
    /**
     * 鐢ㄦ埛涓嬭浇妫€娴嬭褰曪紙鏀寔鏉′欢绛涢€夊悗涓嬭浇锛?
     * @param response HTTP鍝嶅簲瀵硅薄
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     */
    @GetMapping("/download")
    public void downloadUserTestResults(HttpServletResponse response,
                                        @DateTimeFormat(pattern = "yyyy-MM-dd HH:mm:ss") @RequestParam(required = false) LocalDateTime startDate,
                                        @DateTimeFormat(pattern = "yyyy-MM-dd HH:mm:ss") @RequestParam(required = false) LocalDateTime endDate,
                                        @RequestParam(required = false) Boolean isGlandFace,
                                        @RequestParam(required = false) String level) {
        try {
            // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
            Integer userId = UserThreadLocal.getUserId();
            testResultService.exportTestResultsToExcel(response, userId, startDate, endDate, isGlandFace, level);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
    
    /**
     * 鐢ㄦ埛涓嬭浇PDF鏍煎紡妫€娴嬭褰曪紙鏀寔鏉′欢绛涢€夊悗涓嬭浇锛?
     * @param response HTTP鍝嶅簲瀵硅薄
     * @param startDate 寮€濮嬫棩鏈?
     * @param endDate 缁撴潫鏃ユ湡
     * @param isGlandFace 鏄惁涓鸿吅浣撻潰瀹?
     * @param level 绛夌骇
     */
    @GetMapping("/download/pdf")
    public void downloadUserTestResultsAsPDF(HttpServletResponse response,
                                        @DateTimeFormat(pattern = "yyyy-MM-dd HH:mm:ss") @RequestParam(required = false) LocalDateTime startDate,
                                        @DateTimeFormat(pattern = "yyyy-MM-dd HH:mm:ss") @RequestParam(required = false) LocalDateTime endDate,
                                        @RequestParam(required = false) Boolean isGlandFace,
                                        @RequestParam(required = false) String level) {
        try {
            // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
            Integer userId = UserThreadLocal.getUserId();
            testResultService.exportTestResultsToPDF(response, userId, startDate, endDate, isGlandFace, level);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
    
    /**
     * 鐢ㄦ埛鏇存柊鑷繁鐨勬娴嬬粨鏋?
     * @param testResult 妫€娴嬬粨鏋滃璞?
     * @return 鏇存柊缁撴灉
     */
    @PutMapping("/update")
    public Result updateUserTestResult(@RequestBody testResult testResult) {
        try {
            // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
            Integer userId = UserThreadLocal.getUserId();
            boolean success = testResultService.updateUserTestResult(testResult, userId);
            if (success) {
                return Result.success("鏇存柊鎴愬姛");
            } else {
                return Result.error("鏇存柊澶辫触锛岃褰曚笉瀛樺湪鎴栨棤鏉冮檺");
            }
        } catch (Exception e) {
            return Result.error("鏇存柊澶辫触: " + e.getMessage());
        }
    }
}
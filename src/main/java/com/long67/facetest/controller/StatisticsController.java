package com.long67.facetest.controller;

import com.long67.facetest.pojo.Result;
import com.long67.facetest.service.StatisticsService;
import com.long67.facetest.utils.UserThreadLocal;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.time.LocalDate;
import java.util.Map;

@RestController
@RequestMapping("/statistics")
public class StatisticsController {

    @Autowired
    private StatisticsService statisticsService;

    /**
     * 鑾峰彇缁熻姒傝鏁版嵁
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 缁熻姒傝鏁版嵁
     */
    @GetMapping("/overview")
    public Result getStatisticsOverview(
            @DateTimeFormat(pattern = "yyyy-MM-dd") @RequestParam(required = false) LocalDate startDate,
            @DateTimeFormat(pattern = "yyyy-MM-dd") @RequestParam(required = false) LocalDate endDate) {
        try {
            Integer userId = UserThreadLocal.getUserId();
            Map<String, Object> statistics = statisticsService.getStatisticsOverview(userId, startDate, endDate);
            return Result.success(statistics);
        } catch (Exception e) {
            return Result.error("鑾峰彇缁熻姒傝澶辫触: " + e.getMessage());
        }
    }

    /**
     * 鑾峰彇璇︾粏缁熻鏁版嵁
     * @param startDate 寮€濮嬫棩鏈?     * @param endDate 缁撴潫鏃ユ湡
     * @return 璇︾粏缁熻鏁版嵁
     */
    @GetMapping("/detail")
    public Result getStatisticsDetail(
            @DateTimeFormat(pattern = "yyyy-MM-dd") @RequestParam(required = false) LocalDate startDate,
            @DateTimeFormat(pattern = "yyyy-MM-dd") @RequestParam(required = false) LocalDate endDate) {
        try {
            Integer userId = UserThreadLocal.getUserId();
            return Result.success(statisticsService.getStatisticsDetail(userId, startDate, endDate));
        } catch (Exception e) {
            return Result.error("鑾峰彇璇︾粏缁熻澶辫触: " + e.getMessage());
        }
    }
}
package com.long67.facetest.controller;

import com.long67.facetest.pojo.Result;
import com.long67.facetest.pojo.testResult;
import com.long67.facetest.service.DoubaoService;
import com.long67.facetest.utils.UserThreadLocal;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/doubao")
public class DoubaoController {
    
    @Autowired
    private DoubaoService doubaoService;
    
    /**
     * 涓婁紶鍥剧墖骞跺垎鏋愭槸鍚︿负鑵轰綋闈㈠
     * @param image 鍥剧墖鏂囦欢
     * @return 鍒嗘瀽缁撴灉
     */
    @PostMapping("/analyzeFace")
    public Result analyzeFace(@RequestParam("image") MultipartFile image) {
        try {
            // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
            Integer userId = UserThreadLocal.getUserId();
            testResult result = doubaoService.analyzeFace(image, userId);
            return Result.success(result);
        } catch (Exception e) {
            return Result.error("鍒嗘瀽澶辫触: " + e.getMessage());
        }
    }
}
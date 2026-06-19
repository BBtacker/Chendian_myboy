package com.long67.facetest.controller;

import com.long67.facetest.pojo.Result;
import com.long67.facetest.pojo.User;
import com.long67.facetest.service.UserService;
import com.long67.facetest.utils.AliyunOSSOperator;
import com.long67.facetest.utils.UserThreadLocal;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.time.LocalDateTime;

@RestController
@RequestMapping("/user")
public class UserController {
    
    @Autowired
    private UserService userService;
    
    @Autowired
    private AliyunOSSOperator aliyunOSSOperator;
    
    /**
     * 鐢ㄦ埛娉ㄥ唽鏂扮敤鎴?
     */
    @PostMapping("/register")
    public Result<String> register(@RequestBody User user) {
        user.setCreateTime(LocalDateTime.now());
        user.setAvatar(""); // 榛樿澶村儚涓虹┖
        // 璁剧疆榛樿鍊?
        if (!StringUtils.hasText(user.getName())) {
            user.setName(user.getUsername()); // 榛樿濮撳悕涓虹敤鎴峰悕
        }
        userService.add(user);
        return Result.success("娉ㄥ唽鎴愬姛");
    }
    
    /**
     * 鐢ㄦ埛淇敼鑷繁鐨勪俊鎭?
     */
    @PutMapping
    public Result<String> updateUser(@RequestBody User user) {
        // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
        Integer userId = UserThreadLocal.getUserId();
        // 纭繚鐢ㄦ埛鍙兘淇敼鑷繁鐨勪俊鎭?
        user.setId(userId);
        // 濡傛灉瀵嗙爜涓虹┖锛屽垯涓嶆洿鏂板瘑鐮?
        if (!StringUtils.hasText(user.getPassword())) {
            user.setPassword(null);
        }
        userService.update(user);
        return Result.success("鏇存柊鎴愬姛");
    }
    
    /**
     * 鐢ㄦ埛鑾峰彇鑷繁鐨勪俊鎭?
     */
    @GetMapping
    public Result<User> getCurrentUser() {
        // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
        Integer userId = UserThreadLocal.getUserId();
        User user = userService.getById(userId);
        return Result.success(user);
    }
    
    /**
     * 鐢ㄦ埛涓婁紶澶村儚
     */
    @PostMapping("/avatar")
    public Result<String> uploadAvatar(@RequestParam("avatar") MultipartFile avatarFile) {
        try {
            // 浠嶵hreadLocal涓幏鍙栧綋鍓嶇敤鎴稩D
            Integer userId = UserThreadLocal.getUserId();
            
            // 妫€鏌ユ枃浠舵槸鍚︿负绌?
            if (avatarFile.isEmpty()) {
                return Result.error("涓婁紶鐨勫ご鍍忔枃浠朵笉鑳戒负绌?);
            }
            
            // 妫€鏌ユ枃浠剁被鍨?
            String contentType = avatarFile.getContentType();
            if (contentType == null || (!contentType.startsWith("image/jpeg") && 
                                        !contentType.startsWith("image/png") && 
                                        !contentType.startsWith("image/gif"))) {
                return Result.error("鍙敮鎸丣PG銆丳NG銆丟IF鏍煎紡鐨勫浘鐗?);
            }
            
            // 妫€鏌ユ枃浠跺ぇ灏忥紙闄愬埗涓?MB锛?
            if (avatarFile.getSize() > 2 * 1024 * 1024) {
                return Result.error("澶村儚鏂囦欢澶у皬涓嶈兘瓒呰繃2MB");
            }
            
            // 涓婁紶鍒伴樋閲屼簯OSS
            String avatarUrl = aliyunOSSOperator.upload(avatarFile.getBytes(), avatarFile.getOriginalFilename());
            
            // 鏇存柊鐢ㄦ埛澶村儚淇℃伅
            User user = userService.getById(userId);
            user.setAvatar(avatarUrl);
            userService.update(user);
            
            return Result.success(avatarUrl);
        } catch (Exception e) {
            return Result.error("澶村儚涓婁紶澶辫触: " + e.getMessage());
        }
    }
}
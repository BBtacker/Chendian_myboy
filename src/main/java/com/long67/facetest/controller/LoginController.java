package com.long67.facetest.controller;

import com.long67.facetest.pojo.Result;
import com.long67.facetest.pojo.User;
import com.long67.facetest.service.UserService;
import com.long67.facetest.utils.JwtUtils;
import com.long67.facetest.utils.UserThreadLocal;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;

@RestController
public class LoginController {
    
    @Autowired
    private UserService userService;
    
    /**
     * 鐢ㄦ埛鐧诲綍
     */
    @PostMapping("/login")
    public Result login(@RequestBody User user) {
        // 鏍规嵁鐢ㄦ埛鍚嶆煡璇㈢敤鎴?
        User dbUser = userService.getByUsername(user.getUsername());
        
        // 鍒ゆ柇鐢ㄦ埛鏄惁瀛樺湪浠ュ強瀵嗙爜鏄惁姝ｇ‘
        if (dbUser == null) {
            return Result.error("鐢ㄦ埛鍚嶄笉瀛樺湪");
        }
        
        if (!dbUser.getPassword().equals(user.getPassword())) {
            return Result.error("瀵嗙爜閿欒");
        }
        
        // 鐧诲綍鎴愬姛锛岀敓鎴怞WT浠ょ墝
        Map<String, Object> claims = new HashMap<>();
        claims.put("id", dbUser.getId());
        claims.put("username", dbUser.getUsername());
        
        String jwtToken = JwtUtils.genToken(claims);
        
        // 灏嗙敤鎴稩D瀛樺叆ThreadLocal
        UserThreadLocal.setUserId(dbUser.getId());
        
        // 杩斿洖浠ょ墝
        return Result.success(jwtToken);
    }
    
    /**
     * 鐢ㄦ埛鐧诲嚭
     */
    @PostMapping("/logout")
    public Result logout() {
        // 娓呴櫎ThreadLocal涓殑鐢ㄦ埛ID
        UserThreadLocal.clear();
        return Result.success("鐧诲嚭鎴愬姛");
    }
}
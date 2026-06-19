package com.long67.facetest.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.long67.facetest.pojo.Result;
import com.long67.facetest.utils.JwtUtils;
import com.long67.facetest.utils.UserThreadLocal;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.ExpiredJwtException;
import io.jsonwebtoken.MalformedJwtException;
import io.jsonwebtoken.SignatureException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.servlet.HandlerInterceptor;

import java.io.IOException;
import java.util.Date;

@Slf4j
@Component
public class LoginInterceptor implements HandlerInterceptor {
    
    private static final ObjectMapper objectMapper = new ObjectMapper();
    
    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {
        // 鑾峰彇璇锋眰URL
        String url = request.getRequestURL().toString();
        log.info("璇锋眰URL: {}", url);
        
        // 鏀捐鐧诲綍鍜屾敞鍐岀浉鍏虫帴鍙?        if (url.contains("/login") || url.contains("/register")) {
            log.info("鏀捐鐧诲綍/娉ㄥ唽鎺ュ彛");
            return true;
        }
        
        // 浠庤姹傚ご涓幏鍙杢oken
        String token = request.getHeader("Authorization");
        
        // 鍒ゆ柇token鏄惁涓虹┖鎴栨棤鏁?        if (!StringUtils.hasLength(token)) {
            log.info("璇锋眰澶翠腑token涓虹┖");
            responseUnauthorized(response, "鏈櫥褰曪紝璇峰厛鐧诲綍");
            return false;
        }
        
        // 楠岃瘉token
        try {
            Claims claims = JwtUtils.parseJWT(token);
            
            // 妫€鏌oken鏄惁鍗冲皢杩囨湡锛堟彁鍓?0鍒嗛挓鎻愮ず锛?            Date expiration = claims.getExpiration();
            long timeToExpire = expiration.getTime() - System.currentTimeMillis();
            if (timeToExpire < 600000) { // 10鍒嗛挓
                log.info("token鍗冲皢杩囨湡锛屽墿浣欐椂闂? {} 姣", timeToExpire);
            }
            
            // 灏嗙敤鎴蜂俊鎭瓨鍏ヨ姹傚煙涓?            Integer userId = (Integer) claims.get("id");
            String username = (String) claims.get("username");
            
            request.setAttribute("userId", userId);
            request.setAttribute("username", username);
            
            // 灏嗙敤鎴稩D瀛樺叆ThreadLocal涓紝浠ヤ究鍦ㄥ綋鍓嶇嚎绋嬩腑浣跨敤
            UserThreadLocal.setUserId(userId);
            
            log.info("鐢ㄦ埛宸茬櫥褰曪紝鐢ㄦ埛ID: {}, 鐢ㄦ埛鍚? {}", userId, username);
            return true;
        } catch (ExpiredJwtException e) {
            log.info("token宸茶繃鏈? {}", e.getMessage());
            responseUnauthorized(response, "鐧诲綍宸茶繃鏈燂紝璇烽噸鏂扮櫥褰?);
            return false;
        } catch (SignatureException e) {
            log.info("token绛惧悕鏃犳晥: {}", e.getMessage());
            responseUnauthorized(response, "鐧诲綍鍑瘉鏃犳晥锛岃閲嶆柊鐧诲綍");
            return false;
        } catch (MalformedJwtException e) {
            log.info("token鏍煎紡閿欒: {}", e.getMessage());
            responseUnauthorized(response, "鐧诲綍鍑瘉鏍煎紡閿欒锛岃閲嶆柊鐧诲綍");
            return false;
        } catch (Exception e) {
            log.info("token瑙ｆ瀽澶辫触: {}", e.getMessage());
            responseUnauthorized(response, "鐧诲綍鍑瘉楠岃瘉澶辫触锛岃閲嶆柊鐧诲綍");
            return false;
        }
    }
    
    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) throws Exception {
        // 娓呯悊ThreadLocal涓殑鐢ㄦ埛ID锛岄槻姝㈠唴瀛樻硠婕?        UserThreadLocal.clear();
    }
    
    /**
     * 鍝嶅簲鏈巿鏉冮敊璇俊鎭?     */
    private void responseUnauthorized(HttpServletResponse response, String message) throws IOException {
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.setContentType("application/json;charset=utf-8");
        response.getWriter().write(objectMapper.writeValueAsString(Result.error(message)));
    }
}
package com.long67.facetest.utils;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import java.util.Date;
import java.util.Map;

public class JwtUtils {
    private static final String signKey = "faceTest"; // 绛惧悕瀵嗛挜
    private static final Long expire = 86400000L; // 24灏忔椂杩囨湡鏃堕棿(姣)
    
    /**
     * 鐢熸垚JWT浠ょ墝
     * @param claims JWT绗簩閮ㄥ垎璐熻浇(payload)涓瓨鍌ㄧ殑鍐呭
     * @return JWT浠ょ墝
     */
    public static String genToken(Map<String, Object> claims) {
        return Jwts.builder()
                .addClaims(claims) // 鑷畾涔変俊鎭紙杞借嵎锛?                .signWith(SignatureAlgorithm.HS256, signKey) // 绛惧悕绠楁硶鍜屽瘑閽?                .setExpiration(new Date(System.currentTimeMillis() + expire)) // 杩囨湡鏃堕棿
                .compact();
    }
    
    /**
     * 瑙ｆ瀽JWT浠ょ墝
     * @param jwt JWT浠ょ墝
     * @return JWT绗簩閮ㄥ垎璐熻浇(payload)涓瓨鍌ㄧ殑鍐呭
     */
    public static Claims parseJWT(String jwt) {
        return Jwts.parser()
                .setSigningKey(signKey) // 鎸囧畾绛惧悕瀵嗛挜
                .parseClaimsJws(jwt) // 瑙ｆ瀽浠ょ墝
                .getBody();
    }
}
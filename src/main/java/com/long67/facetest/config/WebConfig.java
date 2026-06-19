package com.long67.facetest.config;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class WebConfig implements WebMvcConfigurer {
    
    @Autowired
    private LoginInterceptor loginInterceptor;
    
    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        // 娉ㄥ唽鐧诲綍鎷︽埅鍣紝鎷︽埅鎵€鏈夎姹傦紝鎺掗櫎鐧诲綍鍜屾敞鍐屾帴鍙?        registry.addInterceptor(loginInterceptor)
                .addPathPatterns("/**")  // 鎷︽埅鎵€鏈夎姹?                .excludePathPatterns("/login", "/user/register"); // 鏀捐鐧诲綍鍜屾敞鍐屾帴鍙?    }
}
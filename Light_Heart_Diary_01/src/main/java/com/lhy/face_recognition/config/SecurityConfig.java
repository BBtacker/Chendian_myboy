package com.lhy.face_recognition.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.web.SecurityFilterChain;

/**
 * Spring Security配置类
 */
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
            // 禁用CSRF保护，因为我们使用的是JWT或session-based认证
            .csrf(AbstractHttpConfigurer::disable)
            // 允许所有请求通过，因为我们在Controller中处理认证逻辑
            .authorizeRequests(authorize -> authorize
                .anyRequest().permitAll()
            )
            // 允许跨域请求
            .cors(cors -> cors
                .disable() // 我们在Controller中使用@CrossOrigin注解
            );

        return http.build();
    }
}

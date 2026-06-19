package com.lhy.face_recognition.config;

import com.aliyun.oss.OSSClientBuilder;
import com.aliyun.oss.OSS;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * 阿里云OSS配置类
 */
@Configuration
public class OssConfig {

    @Value("${aliyun.oss.endpoint}")
    private String endpoint;

    @Value("${aliyun.oss.bucketName}")
    private String bucketName;

    @Value("${aliyun.oss.region}")
    private String region;

    @Value("${aliyun.oss.accessKeyId}")
    private String accessKeyId;

    @Value("${aliyun.oss.accessKeySecret}")
    private String accessKeySecret;

    /**
     * 创建OSS客户端实例
     * @return OSS客户端实例
     */
    @Bean
    public OSS ossClient() {
        // 创建OSS客户端实例
        return new OSSClientBuilder().build(endpoint, accessKeyId, accessKeySecret);
    }

    // getter方法，供其他组件使用
    public String getEndpoint() {
        return endpoint;
    }

    public String getBucketName() {
        return bucketName;
    }

    public String getRegion() {
        return region;
    }

    public String getAccessKeyId() {
        return accessKeyId;
    }

    public String getAccessKeySecret() {
        return accessKeySecret;
    }
}

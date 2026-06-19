package com.long67.facetest.utils;

import com.aliyun.oss.*;
import com.aliyun.oss.common.auth.CredentialsProviderFactory;
import com.aliyun.oss.common.auth.DefaultCredentialProvider;
import com.aliyun.oss.common.comm.SignVersion;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.io.ByteArrayInputStream;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.UUID;

@Component
public class AliyunOSSOperator {
    @Autowired
    private AliyunOSSProperties aliyunOSSProperties;
    public String upload(byte[] content, String originalFilename) throws Exception {
        String endpoint = aliyunOSSProperties.getEndpoint();
        String bucketName = aliyunOSSProperties.getBucketName();
        String region = aliyunOSSProperties.getRegion();
        String accessKeyId = aliyunOSSProperties.getAccessKeyId();
        String accessKeySecret = aliyunOSSProperties.getAccessKeySecret();
        
        // 浣跨敤閰嶇疆鏂囦欢涓殑璁块棶鍑瘉
        DefaultCredentialProvider credentialsProvider = new DefaultCredentialProvider(accessKeyId, accessKeySecret);

        // 濉啓Object瀹屾暣璺緞锛屼緥濡?02406/1.png銆侽bject瀹屾暣璺緞涓笉鑳藉寘鍚獴ucket鍚嶇О銆?
        //鑾峰彇褰撳墠绯荤粺鏃ユ湡鐨勫瓧绗︿覆,鏍煎紡涓?yyyy/MM
        String dir = LocalDate.now().format(DateTimeFormatter.ofPattern("yyyy/MM"));
        //鐢熸垚涓€涓柊鐨勪笉閲嶅鐨勬枃浠跺悕
        String newFileName = UUID.randomUUID() + originalFilename.substring(originalFilename.lastIndexOf("."));
        String objectName = dir + "/" + newFileName;

        // 鍒涘缓OSSClient瀹炰緥銆?
        ClientBuilderConfiguration clientBuilderConfiguration = new ClientBuilderConfiguration();
        clientBuilderConfiguration.setSignatureVersion(SignVersion.V4);
        OSS ossClient = OSSClientBuilder.create()
                .endpoint(endpoint)
                .credentialsProvider(credentialsProvider)
                .clientConfiguration(clientBuilderConfiguration)
                .region(region)
                .build();

        try {
            ossClient.putObject(bucketName, objectName, new ByteArrayInputStream(content));
        } finally {
            ossClient.shutdown();
        }

        String protocol = "https://";
        String endpointWithoutProtocol = endpoint.startsWith("http://") || endpoint.startsWith("https://") 
            ? endpoint.substring(endpoint.indexOf("://") + 3) 
            : endpoint;
        return protocol + bucketName + "." + endpointWithoutProtocol + "/" + objectName;
    }

}

package com.lhy.face_recognition.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.RestClientException;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

@Service
public class ExpressionRecognitionService {

    private final RestTemplate restTemplate;

    // Python服务的URL
    private static final String PYTHON_SERVICE_URL = "http://localhost:5000/analyze-expression";

    @Autowired
    public ExpressionRecognitionService(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    /**
     * 调用Python Ultralytics服务进行人脸表情识别
     * @param imageFile 待分析的图像文件
     * @return 表情识别结果
     * @throws IOException 图像读取异常
     */
    public Map<String, Object> recognizeExpression(MultipartFile imageFile) throws IOException {
        // 验证输入文件
        if (imageFile == null || imageFile.isEmpty()) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("error", "未提供有效的图像文件");
            errorResponse.put("total_faces_detected", 0);
            errorResponse.put("faces", new HashMap<>());
            errorResponse.put("expression_counts", new HashMap<>());
            return errorResponse;
        }

        try {
            // 设置HTTP请求头
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.MULTIPART_FORM_DATA);

            // 构建请求体
            MultiValueMap<String, Object> body = new LinkedMultiValueMap<>();
            
            // 使用ByteArrayResource来创建请求体，正确设置文件名和内容类型
            ByteArrayResource byteArrayResource = new ByteArrayResource(imageFile.getBytes()) {
                @Override
                public String getFilename() {
                    return imageFile.getOriginalFilename();
                }
            };
            
            body.add("image", byteArrayResource);

            // 创建HTTP请求实体
            HttpEntity<MultiValueMap<String, Object>> requestEntity = new HttpEntity<>(body, headers);

            // 调用Python服务
            ResponseEntity<Map> response = restTemplate.postForEntity(PYTHON_SERVICE_URL, requestEntity, Map.class);
            return response.getBody();
        } catch (HttpServerErrorException e) {
            // 捕获Python服务返回的500等错误
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("error", "表情识别服务处理失败: " + e.getResponseBodyAsString());
            errorResponse.put("total_faces_detected", 0);
            errorResponse.put("faces", new HashMap<>());
            errorResponse.put("expression_counts", new HashMap<>());
            return errorResponse;
        } catch (RestClientException e) {
            // 捕获其他REST客户端异常
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("error", "与表情识别服务通信失败: " + e.getMessage());
            errorResponse.put("total_faces_detected", 0);
            errorResponse.put("faces", new HashMap<>());
            errorResponse.put("expression_counts", new HashMap<>());
            return errorResponse;
        } catch (Exception e) {
            // 捕获所有其他异常
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("error", "表情识别过程中发生未知错误: " + e.getMessage());
            errorResponse.put("total_faces_detected", 0);
            errorResponse.put("faces", new HashMap<>());
            errorResponse.put("expression_counts", new HashMap<>());
            return errorResponse;
        }
    }
}
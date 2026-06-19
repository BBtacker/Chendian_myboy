package com.lhy.face_recognition.service.impl;

import com.aliyun.oss.OSS;
import com.aliyun.oss.OSSClientBuilder;
import com.aliyun.oss.model.PutObjectRequest;
import com.lhy.face_recognition.service.DoubaoService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.io.ByteArrayInputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@Service
public class DoubaoServiceImpl implements DoubaoService {

    @Autowired
    private RestTemplate restTemplate;

    // 从配置文件获取API密钥
    @Value("${doubao.api.key}")
    private String apiKey;

    // OSS配置
    @Value("${aliyun.oss.endpoint}")
    private String ossEndpoint;
    @Value("${aliyun.oss.bucketName}")
    private String ossBucketName;
    @Value("${aliyun.oss.accessKeyId}")
    private String ossAccessKeyId;
    @Value("${aliyun.oss.accessKeySecret}")
    private String ossAccessKeySecret;

    // 模型ID
    private static final String MODEL_ID = "doubao-seed-2-0-pro-260215";
    // 图片生成模型ID
    private static final String IMAGE_MODEL_ID = "doubao-seedream-4-0-250828";

    // 请求地址（以北京地域为例）
    private static final String API_URL = "https://ark.cn-beijing.volces.com/api/v3/chat/completions";
    // 图片生成API地址
    private static final String IMAGE_API_URL = "https://ark.cn-beijing.volces.com/api/v3/images/generations";

    // 优化的系统提示词，更符合应用场景
    private static final String SYSTEM_PROMPT = "你是一个专注于日记、心情或生活记录的智能助手，名叫心影助手。请使用温暖、亲切的语气，帮助用户解决关于日记写作、心情管理和生活记录的问题。无论用户问什么问题，都尽量与日记和心情扯上关系，提供具体、有建设性的回答，并且富有同理心。回答要自然流畅，不要生硬地转折话题，而是巧妙地将任何问题与日记或心情管理联系起来。";

    /**
     * 将 base64 图片上传到阿里云 OSS
     * @param base64Data base64 图片数据
     * @return OSS 上的图片 URL
     */
    private String uploadImageToOSS(String base64Data) throws Exception {
        System.out.println("=== 开始上传图片到 OSS ===");
        System.out.println("OSS Endpoint: " + ossEndpoint);
        System.out.println("OSS Bucket: " + ossBucketName);
        System.out.println("OSS AccessKeyId: " + (ossAccessKeyId != null ? ossAccessKeyId.substring(0, 10) + "..." : "null"));
        
        // 移除 base64 前缀
        String pureBase64 = base64Data;
        if (base64Data.contains(",")) {
            pureBase64 = base64Data.split(",")[1];
            System.out.println("已移除 base64 前缀");
        }
        
        // 解码 base64 数据
        byte[] imageBytes = java.util.Base64.getDecoder().decode(pureBase64);
        System.out.println("图片字节大小: " + imageBytes.length);
        
        // 生成唯一的文件名
        String fileName = "dream-images/" + UUID.randomUUID().toString() + ".png";
        System.out.println("OSS 文件名: " + fileName);
        
        // 创建 OSS 客户端
        OSS ossClient = new OSSClientBuilder().build(ossEndpoint, ossAccessKeyId, ossAccessKeySecret);
        
        try {
            // 上传图片
            PutObjectRequest putObjectRequest = new PutObjectRequest(ossBucketName, fileName, new ByteArrayInputStream(imageBytes));
            ossClient.putObject(putObjectRequest);
            System.out.println("图片上传成功");
            
            // 生成 URL
            String url = "https://" + ossBucketName + "." + ossEndpoint + "/" + fileName;
            System.out.println("生成的图片 URL: " + url);
            return url;
        } finally {
            // 关闭 OSS 客户端
            ossClient.shutdown();
            System.out.println("OSS 客户端已关闭");
        }
    }

    @Override
    public String getDoubaoAnswer(List<Map<String, Object>> messages) {
        try {
            // 确保messages格式正确
            List<Map<String, Object>> formattedMessages = new ArrayList<>();
            
            // 添加系统提示词
            Map<String, Object> systemMessage = new HashMap<>();
            systemMessage.put("role", "system");
            systemMessage.put("content", SYSTEM_PROMPT);
            formattedMessages.add(systemMessage);
            
            // 处理用户和助手消息
            for (Map<String, Object> msg : messages) {
                Map<String, Object> formattedMsg = new HashMap<>();
                // 将Object类型转换为boolean类型
                boolean isUser = msg.get("isUser") != null && Boolean.parseBoolean(msg.get("isUser").toString());
                formattedMsg.put("role", isUser ? "user" : "assistant");
                formattedMsg.put("content", msg.get("content"));
                formattedMessages.add(formattedMsg);
            }

            // 请求头
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            headers.set("Authorization", "Bearer " + apiKey);

            // 请求参数，添加温度和最大令牌数，优化回答效果
            Map<String, Object> requestBody = new HashMap<>();
            requestBody.put("model", MODEL_ID);
            requestBody.put("messages", formattedMessages);
            requestBody.put("temperature", 0.7);
            requestBody.put("max_tokens", 1000);
            requestBody.put("top_p", 0.9);

            // 发起请求
            HttpEntity<Map<String, Object>> requestEntity = new HttpEntity<>(requestBody, headers);
            ResponseEntity<Map> response = restTemplate.postForEntity(API_URL, requestEntity, Map.class);

            // 处理响应结果
            Map<String, Object> result = response.getBody();
            if (result != null && result.containsKey("choices")) {
                List<Map<String, Object>> choices = (List<Map<String, Object>>) result.get("choices");
                if (!choices.isEmpty()) {
                    Map<String, Object> choice = choices.get(0);
                    Map<String, Object> message = (Map<String, Object>) choice.get("message");
                    return (String) message.get("content");
                }
            }

            return "抱歉，我暂时无法回答您的问题。";
        } catch (Exception e) {
            e.printStackTrace();
            return "抱歉，我暂时无法回答您的问题，请稍后再试。";
        }
    }

    @Override
    public String generateImage(String imageData, String mood) {
        System.out.println("=== 开始生成图片 ===");
        System.out.println("Mood: " + mood);
        System.out.println("Image data length: " + (imageData != null ? imageData.length() : 0));
        System.out.println("API Key: " + (apiKey != null ? apiKey.substring(0, 10) + "..." : "null"));
        System.out.println("Image API URL: " + IMAGE_API_URL);
        
        try {
            // 1. 上传图片到 OSS 获取 URL
            String imageUrl = uploadImageToOSS(imageData);
            System.out.println("获取到图片 URL: " + imageUrl);
            
            // 2. 构建提示词
            String prompt = String.format(
                "根据上述图片，生成完整保留原图的主体内容与画面构图，整体画面干净通透，带着柔和不刺眼的自然光线，用清新淡雅的暖色调铺色，线条干净利落，整体氛围安静又温馨，充满了松弛治愈的感觉，画质细腻精致，达到动画原画级别的顶级水准。心情类型：%s",
                mood
            );
            System.out.println("Prompt: " + prompt);

            // 3. 构建请求体 - 只包含模型支持的参数
            String requestBody = String.format(
                "{\"model\":\"%s\",\"prompt\":\"%s\",\"image\":\"%s\",\"response_format\":\"url\",\"watermark\":false}",
                IMAGE_MODEL_ID, prompt, imageUrl
            );
            System.out.println("Request body length: " + requestBody.length());
            System.out.println("Request body: " + requestBody);

            // 4. 请求头
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            headers.set("Authorization", "Bearer " + apiKey);
            System.out.println("Headers: " + headers);

            // 5. 发起请求
            HttpEntity<String> requestEntity = new HttpEntity<>(requestBody, headers);
            System.out.println("发起图片生成请求...");
            
            long startTime = System.currentTimeMillis();
            ResponseEntity<String> response = restTemplate.postForEntity(
                IMAGE_API_URL, requestEntity, String.class
            );
            long endTime = System.currentTimeMillis();
            
            System.out.println("请求耗时: " + (endTime - startTime) + "ms");
            System.out.println("图片生成请求响应状态码: " + response.getStatusCode());
            System.out.println("响应头: " + response.getHeaders());
            
            String responseBody = response.getBody();
            System.out.println("图片生成请求响应结果长度: " + (responseBody != null ? responseBody.length() : 0));
            System.out.println("图片生成请求响应结果: " + responseBody);
            
            // 6. 检查响应
            if (responseBody != null && !responseBody.isEmpty()) {
                // 解析响应 JSON，提取图片 URL
                System.out.println("=== 图片生成成功 ===");
                try {
                    // 使用简单的字符串处理提取 URL
                    if (responseBody.contains("url")) {
                        int urlStart = responseBody.indexOf("url\":\"") + 6;
                        int urlEnd = responseBody.indexOf("\"", urlStart);
                        if (urlStart > 0 && urlEnd > urlStart) {
                            String imageUrlResult = responseBody.substring(urlStart, urlEnd);
                            System.out.println("提取到图片 URL: " + imageUrlResult);
                            return imageUrlResult;
                        }
                    }
                    System.out.println("无法从响应中提取 URL，返回原始响应");
                    return responseBody;
                } catch (Exception e) {
                    System.out.println("解析响应失败: " + e.getMessage());
                    return responseBody;
                }
            } else {
                System.out.println("=== 响应为空 ===");
                return null;
            }
        } catch (Exception e) {
            System.out.println("=== 生成图片异常 ===");
            e.printStackTrace();
            return null;
        }
    }
}
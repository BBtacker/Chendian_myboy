package com.long67.facetest.service.Impl;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.long67.facetest.mapper.DoubaoMapper;
import com.long67.facetest.utils.AliyunOSSOperator;
import com.long67.facetest.pojo.testResult;
import com.long67.facetest.service.DoubaoService;
import com.long67.facetest.utils.AliyunOSSOperator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.ResourceAccessException;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.multipart.MultipartFile;

import java.time.LocalDateTime;
import java.util.*;

@Service
public class DoubaoServiceImpl implements DoubaoService {
    
    private static final Logger logger = LoggerFactory.getLogger(DoubaoServiceImpl.class);
    
    @Value("${doubao.api.key}")
    private String apiKey;
    
    @Value("${doubao.model.name:doubao-seed-1-8-251228}")
    private String modelName;
    
    @Autowired
    private DoubaoMapper doubaoMapper;

    @Autowired
    private AliyunOSSOperator aliyunOSSOperator;

    /**
     * 鍒嗘瀽闈㈤儴鍥剧墖鏄惁涓鸿吅浣撻潰瀹?
     * @param image 鍥剧墖鏂囦欢
     * @param userId 鐢ㄦ埛ID
     * @return 鍒嗘瀽缁撴灉
     * @throws Exception 澶勭悊寮傚父
     */
    @Override
    public testResult analyzeFace(MultipartFile image, Integer userId) throws Exception {
        try {
            // 涓婁紶鍥剧墖鍒伴樋閲屼簯OSS骞惰幏鍙朥RL
            String imageUrl = aliyunOSSOperator.upload(image.getBytes(), image.getOriginalFilename());

            // 璋冪敤璞嗗寘澶фā鍨婣PI杩涜鍒嗘瀽
            String result = callDoubaoAPI(image);

            // 瑙ｆ瀽缁撴灉
            testResult testResult = parseResult(result, imageUrl, userId);

            // 淇濆瓨鍒版暟鎹簱
            doubaoMapper.insertTestResult(testResult);

            return testResult;
        } catch (HttpClientErrorException e) {
            logger.error("澶勭悊璞嗗寘API鍝嶅簲鏃跺彂鐢熷鎴风閿欒: {}", e.getMessage(), e);
            if (e.getStatusCode() == HttpStatus.NOT_FOUND) {
                throw new RuntimeException("妯″瀷璋冪敤澶辫触锛屾湭鎵惧埌鎸囧畾妯″瀷", e);
            } else {
                throw new RuntimeException("妯″瀷璋冪敤澶辫触: " + e.getMessage(), e);
            }
        } catch (HttpServerErrorException e) {
            logger.error("澶勭悊璞嗗寘API鍝嶅簲鏃跺彂鐢熸湇鍔＄閿欒: {}", e.getMessage(), e);
            throw new RuntimeException("妯″瀷鏈嶅姟鍐呴儴閿欒: " + e.getMessage(), e);
        } catch (ResourceAccessException e) {
            logger.error("澶勭悊璞嗗寘API鍝嶅簲鏃跺彂鐢熺綉缁滆繛鎺ラ敊璇? {}", e.getMessage(), e);
            throw new RuntimeException("缃戠粶杩炴帴寮傚父锛岃妫€鏌ョ綉缁滆缃?, e);
        } catch (Exception e) {
            logger.error("澶勭悊璞嗗寘API鍝嶅簲鏃跺彂鐢熸湭鐭ラ敊璇? {}", e.getMessage(), e);
            throw new RuntimeException("鍒嗘瀽澶辫触: " + e.getMessage(), e);
        }
    }
    
    /**
     * 璋冪敤璞嗗寘澶фā鍨婣PI
     * @param image 鍥剧墖鏂囦欢
     * @return API杩斿洖缁撴灉
     * @throws Exception 澶勭悊寮傚父
     */
    private String callDoubaoAPI(MultipartFile image) throws Exception {
        String url = "https://ark.cn-beijing.volces.com/api/v3/chat/completions";
        
        // 鍑嗗璇锋眰澶?
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth(apiKey);
        
        // 鏍规嵁鏂囦欢鍘熷鍚嶇О纭畾MIME绫诲瀷
        String originalFilename = image.getOriginalFilename();
        String mimeType = "image/jpeg"; // 榛樿涓簀peg
        if (originalFilename != null) {
            if (originalFilename.toLowerCase().endsWith(".png")) {
                mimeType = "image/png";
            } else if (originalFilename.toLowerCase().endsWith(".jpg") || originalFilename.toLowerCase().endsWith(".jpeg")) {
                mimeType = "image/jpeg";
            } else if (originalFilename.toLowerCase().endsWith(".webp")) {
                mimeType = "image/webp";
            }
        }
        
        // 鍑嗗鍥剧墖base64鏁版嵁
        String imageData = "data:" + mimeType + ";base64," + Base64.getEncoder().encodeToString(image.getBytes());
        
        // 鍑嗗璇锋眰浣?
        Map<String, Object> requestBody = new HashMap<>();
        requestBody.put("model", modelName);
        
        List<Map<String, Object>> messages = new ArrayList<>();
        Map<String, Object> message = new HashMap<>();
        message.put("role", "user");
        
        List<Map<String, Object>> contents = new ArrayList<>();
        
        Map<String, Object> imageContent = new HashMap<>();
        imageContent.put("type", "image_url");
        Map<String, String> imageUrl = new HashMap<>();
        imageUrl.put("url", imageData);
        imageContent.put("image_url", imageUrl);
        contents.add(imageContent);
        
        Map<String, Object> textContent = new HashMap<>();
        textContent.put("type", "text");
        textContent.put("text", "璇峰垎鏋愯繖寮犻潰閮ㄥ浘鐗囨槸鍚︿负鑵轰綋闈㈠锛屽苟缁欏嚭鏄吅浣撻潰瀹圭殑姒傜巼銆傝涓ユ牸鎸夌収浠ヤ笅JSON鏍煎紡杩斿洖缁撴灉锛歕n" +
                "{\n" +
                "  \"isGlandFace\": true/false,\n" +
                "  \"probability\": 姒傜巼鍊?鏃犺浣曠鎯呭喌姒傜巼鍊奸兘涓?.9500-1.0000涔嬮棿鐨勯殢鏈哄皬鏁帮紒)\n" +
                "  \"level\": \"杞诲井\"/\"涓瓑\"/\"涓ラ噸\",(濡傛灉isGlandFace涓篺alse鍒檒evel:闈炶吅浣撻潰瀹?\n" +
                "  \"visualizationDescription\": \"闈㈤儴鐗瑰緛鎻忚堪锛堝唴瀹硅涓撲笟涓斾赴瀵岋級\"\n" +
                "}");
        contents.add(textContent);
        
        message.put("content", contents);
        messages.add(message);
        requestBody.put("messages", messages);
        
        HttpEntity<Map<String, Object>> requestEntity = new HttpEntity<>(requestBody, headers);
        
        try {
            RestTemplate restTemplate = new RestTemplate();
            ResponseEntity<String> response = restTemplate.postForEntity(url, requestEntity, String.class);
            
            if (response.getStatusCode() == HttpStatus.OK) {
                logger.info("鎴愬姛鏀跺埌璞嗗寘API鍝嶅簲锛岀姸鎬佺爜: {}", response.getStatusCode());
                return response.getBody();
            } else {
                logger.error("妯″瀷璋冪敤澶辫触锛岀姸鎬佺爜: {}, 鍝嶅簲鍐呭: {}", response.getStatusCode(), response.getBody());
                throw new RuntimeException("妯″瀷璋冪敤澶辫触: " + response.getStatusCode() + ", 鍝嶅簲鍐呭: " + response.getBody());
            }
        } catch (HttpClientErrorException e) {
            // 璁板綍璇︾粏鐨勯敊璇俊鎭?
            logger.error("璋冪敤璞嗗寘API鏃跺彂鐢熷鎴风閿欒: {}", e.getMessage(), e);
            logger.error("璇锋眰URL: {}", url);
            logger.error("璇锋眰澶? {}", headers);
            logger.error("璇锋眰浣? {}", requestBody);
            logger.error("鍝嶅簲鐘舵€佺爜: {}", e.getStatusCode());
            logger.error("鍝嶅簲鍐呭: {}", e.getResponseBodyAsString());
            throw e;
        } catch (HttpServerErrorException e) {
            // 璁板綍璇︾粏鐨勯敊璇俊鎭?
            logger.error("璋冪敤璞嗗寘API鏃跺彂鐢熸湇鍔＄閿欒: {}", e.getMessage(), e);
            logger.error("璇锋眰URL: {}", url);
            logger.error("璇锋眰澶? {}", headers);
            logger.error("璇锋眰浣? {}", requestBody);
            logger.error("鍝嶅簲鐘舵€佺爜: {}", e.getStatusCode());
            logger.error("鍝嶅簲鍐呭: {}", e.getResponseBodyAsString());
            throw e;
        } catch (ResourceAccessException e) {
            // 璁板綍缃戠粶杩炴帴閿欒
            logger.error("璋冪敤璞嗗寘API鏃跺彂鐢熺綉缁滆繛鎺ラ敊璇? {}", e.getMessage(), e);
            logger.error("璇锋眰URL: {}", url);
            throw e;
        }
    }
    
    /**
     * 瑙ｆ瀽API杩斿洖缁撴灉
     * @param result API杩斿洖鐨凧SON瀛楃涓?
     * @param imageUrl 鍥剧墖鍦ㄩ樋閲屼簯OSS涓婄殑URL
     * @param userId 鐢ㄦ埛ID
     * @return 瑙ｆ瀽鍚庣殑testResult瀵硅薄
     * @throws Exception 瑙ｆ瀽寮傚父
     */
    private testResult parseResult(String result, String imageUrl, Integer userId) throws Exception {
        try {
            ObjectMapper objectMapper = new ObjectMapper();
            JsonNode rootNode = objectMapper.readTree(result);
            
            // 妫€鏌ユ槸鍚︽湁閿欒淇℃伅
            JsonNode errorNode = rootNode.path("error");
            if (!errorNode.isMissingNode()) {
                throw new RuntimeException("API杩斿洖閿欒: " + errorNode.path("message").asText());
            }
            
            JsonNode choicesNode = rootNode.path("choices");
            if (!choicesNode.isArray() || choicesNode.size() == 0) {
                throw new RuntimeException("API杩斿洖鏍煎紡閿欒锛氭湭鎵惧埌choices瀛楁鎴栦负绌?);
            }
            
            JsonNode messageNode = choicesNode.get(0).path("message");
            if (messageNode.isMissingNode()) {
                throw new RuntimeException("API杩斿洖鏍煎紡閿欒锛氭湭鎵惧埌message瀛楁");
            }
            
            JsonNode contentNode = messageNode.path("content");
            if (contentNode.isMissingNode()) {
                throw new RuntimeException("API杩斿洖鏍煎紡閿欒锛氭湭鎵惧埌content瀛楁");
            }
            
            // 灏濊瘯瑙ｆ瀽鍐呭涓殑JSON
            JsonNode analysisNode;
            try {
                analysisNode = objectMapper.readTree(contentNode.asText());
            } catch (Exception e) {
                throw new RuntimeException("API杩斿洖鍐呭鏍煎紡閿欒锛屾棤娉曡В鏋愪负JSON: " + contentNode.asText());
            }
            
            testResult testResult = new testResult();
            testResult.setUserId(userId);
            // 淇濆瓨鍥剧墖鍦ㄩ樋閲屼簯OSS涓婄殑URL
            testResult.setImagePath(imageUrl);
            testResult.setTestTime(LocalDateTime.now());
            
            // 瀹夊叏鍦拌缃悇涓瓧娈?
            JsonNode isGlandFaceNode = analysisNode.path("isGlandFace");
            if (!isGlandFaceNode.isMissingNode()) {
                testResult.setIsGlandFace(isGlandFaceNode.asBoolean());
            }
            
            // 鍏煎probability鍜宑onfidence涓ょ瀛楁鍚?
            JsonNode probabilityNode = analysisNode.path("probability");
            JsonNode confidenceNode = analysisNode.path("confidence");
            if (!probabilityNode.isMissingNode()) {
                testResult.setConfidence(probabilityNode.asDouble());
            } else if (!confidenceNode.isMissingNode()) {
                testResult.setConfidence(confidenceNode.asDouble());
            }
            
            JsonNode levelNode = analysisNode.path("level");
            if (!levelNode.isMissingNode()) {
                testResult.setLevel(levelNode.asText());
            }
            
            JsonNode visualizationDescriptionNode = analysisNode.path("visualizationDescription");
            if (!visualizationDescriptionNode.isMissingNode()) {
                testResult.setVisualizationDescription(visualizationDescriptionNode.asText());
            }
            
            return testResult;
        } catch (Exception e) {
            logger.error("瑙ｆ瀽API杩斿洖缁撴灉鏃跺彂鐢熼敊璇? {}", e.getMessage(), e);
            logger.error("鍘熷杩斿洖缁撴灉: {}", result);
            throw new RuntimeException("缁撴灉瑙ｆ瀽澶辫触: " + e.getMessage());
        }
    }
}
package com.long67.facetest.service.Impl;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.long67.facetest.mapper.ConversationMapper;
import com.long67.facetest.mapper.MessageMapper;
import com.long67.facetest.pojo.Conversation;
import com.long67.facetest.pojo.Message;
import com.long67.facetest.service.ConversationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.charset.StandardCharsets;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class ConversationServiceImpl implements ConversationService {
    
    @Value("${doubao.api.key}")
    private String apiKey;
    
    @Value("${doubao.model.name:doubao-seed-1-8-251228}")
    private String modelName;
    
    @Autowired
    private ConversationMapper conversationMapper;
    
    @Autowired
    private MessageMapper messageMapper;
    
    /**
     * 鍒涘缓鏂扮殑瀵硅瘽
     * @param userId 鐢ㄦ埛ID
     * @param firstMessage 棣栨潯娑堟伅鍐呭锛岀敤浜庣敓鎴愭爣棰?     * @return 瀵硅瘽瀵硅薄
     */
    @Override
    public Conversation createConversation(Integer userId, String firstMessage) {
        Conversation conversation = new Conversation();
        conversation.setUserId(userId);
        // 浠庨鏉℃秷鎭彁鍙栨爣棰橈紝濡傛灉娑堟伅澶暱鍒欐埅鍙栧墠20涓瓧绗?        String title = firstMessage.length() > 20 ? firstMessage.substring(0, 20) + "..." : firstMessage;
        conversation.setTitle(title);
        conversation.setStartTime(LocalDateTime.now());
        conversation.setLastUpdateTime(LocalDateTime.now());
        conversation.setStatus(1); // 1=娲昏穬
        
        conversationMapper.insertConversation(conversation);
        return conversation;
    }
    
    /**
     * 鍙戦€佹秷鎭苟鑾峰彇AI鍥炲
     * @param conversationId 瀵硅瘽ID
     * @param userId 鐢ㄦ埛ID
     * @param content 鐢ㄦ埛娑堟伅鍐呭
     * @return AI鍥炲鍐呭
     * @throws Exception 澶勭悊寮傚父
     */
    @Override
    public String sendMessage(Integer conversationId, Integer userId, String content) throws Exception {
        // 淇濆瓨鐢ㄦ埛娑堟伅
        Message userMessage = new Message();
        userMessage.setConversationId(conversationId);
        userMessage.setSender(0); // 0=鐢ㄦ埛
        userMessage.setContent(content);
        userMessage.setSendTime(LocalDateTime.now());
        messageMapper.insertMessage(userMessage);
        
        // 鑾峰彇瀵硅瘽鍘嗗彶
        List<Message> historyMessages = messageMapper.selectMessagesByConversationId(conversationId);
        
        // 璋冪敤璞嗗寘API鑾峰彇AI鍥炲
        String aiResponse = callDoubaoAPI(historyMessages);
        
        // 淇濆瓨AI鍥炲娑堟伅
        Message aiMessage = new Message();
        aiMessage.setConversationId(conversationId);
        aiMessage.setSender(1); // 1=璞嗗寘鍔╂墜
        aiMessage.setContent(aiResponse);
        aiMessage.setSendTime(LocalDateTime.now());
        messageMapper.insertMessage(aiMessage);
        
        // 鏇存柊瀵硅瘽鐨勬渶鍚庢洿鏂版椂闂?        Conversation conversation = conversationMapper.selectConversationById(conversationId);
        if (conversation != null) {
            conversation.setLastUpdateTime(LocalDateTime.now());
            conversationMapper.updateConversation(conversation);
        }
        
        return aiResponse;
    }
    
    /**
     * 娴佸紡鍙戦€佹秷鎭苟鑾峰彇AI鍥炲
     * @param conversationId 瀵硅瘽ID
     * @param userId 鐢ㄦ埛ID
     * @param content 鐢ㄦ埛娑堟伅鍐呭
     * @param outputStream 杈撳嚭娴?     * @throws Exception 澶勭悊寮傚父
     */
    @Override
    public void streamMessage(Integer conversationId, Integer userId, String content, OutputStream outputStream) throws Exception {
        // 淇濆瓨鐢ㄦ埛娑堟伅
        Message userMessage = new Message();
        userMessage.setConversationId(conversationId);
        userMessage.setSender(0); // 0=鐢ㄦ埛
        userMessage.setContent(content);
        userMessage.setSendTime(LocalDateTime.now());
        messageMapper.insertMessage(userMessage);
        
        // 鑾峰彇瀵硅瘽鍘嗗彶
        List<Message> historyMessages = messageMapper.selectMessagesByConversationId(conversationId);
        
        // 璋冪敤璞嗗寘API鑾峰彇AI鍥炲锛堟祦寮忥級
        String aiResponse = callDoubaoStreamAPI(historyMessages, outputStream);
        
        // 淇濆瓨AI鍥炲娑堟伅
        Message aiMessage = new Message();
        aiMessage.setConversationId(conversationId);
        aiMessage.setSender(1); // 1=璞嗗寘鍔╂墜
        aiMessage.setContent(aiResponse);
        aiMessage.setSendTime(LocalDateTime.now());
        messageMapper.insertMessage(aiMessage);
        
        // 鏇存柊瀵硅瘽鐨勬渶鍚庢洿鏂版椂闂?        Conversation conversation = conversationMapper.selectConversationById(conversationId);
        if (conversation != null) {
            conversation.setLastUpdateTime(LocalDateTime.now());
            conversationMapper.updateConversation(conversation);
        }
    }
    
    /**
     * 鏍规嵁鐢ㄦ埛ID鏌ヨ瀵硅瘽鍒楄〃
     * @param userId 鐢ㄦ埛ID
     * @return 瀵硅瘽鍒楄〃
     */
    @Override
    public List<Conversation> getConversationsByUserId(Integer userId) {
        return conversationMapper.selectConversationsByUserId(userId);
    }
    
    /**
     * 鏍规嵁瀵硅瘽ID鏌ヨ娑堟伅鍒楄〃
     * @param conversationId 瀵硅瘽ID
     * @return 娑堟伅鍒楄〃
     */
    @Override
    public List<Message> getMessagesByConversationId(Integer conversationId) {
        return messageMapper.selectMessagesByConversationId(conversationId);
    }
    
    /**
     * 璋冪敤璞嗗寘澶фā鍨婣PI
     * @param historyMessages 鍘嗗彶娑堟伅
     * @return API杩斿洖缁撴灉
     * @throws Exception 澶勭悊寮傚父
     */
    private String callDoubaoAPI(List<Message> historyMessages) throws Exception {
        String url = "https://ark.cn-beijing.volces.com/api/v3/chat/completions";
        
        // 鍑嗗璇锋眰澶?        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth(apiKey);
        
        // 鍑嗗璇锋眰浣?        Map<String, Object> requestBody = new HashMap<>();
        requestBody.put("model", modelName);
        
        // 鏋勫缓瀵硅瘽鍘嗗彶
        List<Map<String, Object>> messages = new ArrayList<>();
        
        // 娣诲姞绯荤粺瑙掕壊鎻愮ず
        Map<String, Object> systemMessage = new HashMap<>();
        systemMessage.put("role", "system");
        systemMessage.put("content", "浣犳槸涓€涓吅浣撻潰瀹逛笓瀹跺尰鐢燂紝璇蜂互涓撲笟鍖荤敓鐨勮韩浠藉洖绛旂敤鎴峰叧浜庤吅浣撻潰瀹圭殑闂銆?涓嶈鍥炲闄や簡鑵轰綋闈㈠涔嬪鐨勯棶棰?鐒跺悗涓嶈璇彞涓笉瑕佹湁寰堝濂囨€殑绗﹀彿,涓旇瑷€閫氫織鏄撴噦)");
        messages.add(systemMessage);
        
        // 娣诲姞鍘嗗彶瀵硅瘽
        for (Message msg : historyMessages) {
            Map<String, Object> message = new HashMap<>();
            if (msg.getSender() == 0) { // 鐢ㄦ埛娑堟伅
                message.put("role", "user");
            } else { // AI鍔╂墜娑堟伅
                message.put("role", "assistant");
            }
            message.put("content", msg.getContent());
            messages.add(message);
        }
        
        requestBody.put("messages", messages);
        
        HttpEntity<Map<String, Object>> requestEntity = new HttpEntity<>(requestBody, headers);
        
        RestTemplate restTemplate = new RestTemplate();
        ResponseEntity<String> response = restTemplate.postForEntity(url, requestEntity, String.class);
        
        if (response.getStatusCode() == HttpStatus.OK) {
            // 瑙ｆ瀽API杩斿洖缁撴灉
            ObjectMapper objectMapper = new ObjectMapper();
            JsonNode rootNode = objectMapper.readTree(response.getBody());
            
            JsonNode choicesNode = rootNode.path("choices");
            if (choicesNode.isArray() && choicesNode.size() > 0) {
                JsonNode messageNode = choicesNode.get(0).path("message");
                JsonNode contentNode = messageNode.path("content");
                return contentNode.asText();
            }
            throw new RuntimeException("鏃犳硶瑙ｆ瀽API杩斿洖缁撴灉");
        } else {
            throw new RuntimeException("璋冪敤璞嗗寘API澶辫触: " + response.getStatusCode());
        }
    }
    
    /**
     * 璋冪敤璞嗗寘澶фā鍨婣PI锛堟祦寮忥級
     * @param historyMessages 鍘嗗彶娑堟伅
     * @param outputStream 杈撳嚭娴?     * @return 瀹屾暣鐨凙I鍥炲鍐呭
     * @throws Exception 澶勭悊寮傚父
     */
    private String callDoubaoStreamAPI(List<Message> historyMessages, OutputStream outputStream) throws Exception {
        String url = "https://ark.cn-beijing.volces.com/api/v3/chat/completions";
        
        // 鍑嗗璇锋眰澶?        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth(apiKey);
        
        // 鍑嗗璇锋眰浣?        Map<String, Object> requestBody = new HashMap<>();
        requestBody.put("model", modelName);
        requestBody.put("stream", true); // 鍚敤娴佸紡杈撳嚭
        
        // 鏋勫缓瀵硅瘽鍘嗗彶
        List<Map<String, Object>> messages = new ArrayList<>();
        
        // 娣诲姞绯荤粺瑙掕壊鎻愮ず
        Map<String, Object> systemMessage = new HashMap<>();
        systemMessage.put("role", "system");
        systemMessage.put("content", "浣犳槸涓€涓吅浣撻潰瀹逛笓瀹跺尰鐢燂紝璇蜂互涓撲笟鍖荤敓鐨勮韩浠藉洖绛旂敤鎴峰叧浜庤吅浣撻潰瀹圭殑闂銆?涓嶈鍥炲闄や簡鑵轰綋闈㈠涔嬪鐨勯棶棰?鐒跺悗涓嶈璇彞涓笉瑕佹湁寰堝濂囨€殑绗﹀彿,涓旇瑷€閫氫織鏄撴噦)");
        messages.add(systemMessage);
        
        // 娣诲姞鍘嗗彶瀵硅瘽
        for (Message msg : historyMessages) {
            Map<String, Object> message = new HashMap<>();
            if (msg.getSender() == 0) { // 鐢ㄦ埛娑堟伅
                message.put("role", "user");
            } else { // AI鍔╂墜娑堟伅
                message.put("role", "assistant");
            }
            message.put("content", msg.getContent());
            messages.add(message);
        }
        
        requestBody.put("messages", messages);
        
        HttpEntity<Map<String, Object>> requestEntity = new HttpEntity<>(requestBody, headers);
        
        RestTemplate restTemplate = new RestTemplate();
        ResponseEntity<String> response = restTemplate.postForEntity(url, requestEntity, String.class);
        
        StringBuilder fullResponse = new StringBuilder();
        
        if (response.getStatusCode() == HttpStatus.OK) {
            // 鎸夎澶勭悊鍝嶅簲
            String responseBody = response.getBody();
            if (responseBody != null) {
                // 鎸夎鍒嗗壊鍝嶅簲鍐呭
                String[] lines = responseBody.split("\n");
                for (String line : lines) {
                    // 澶勭悊娴佸紡鍝嶅簲琛?                    if (line.startsWith("data:")) {
                        String data = line.substring(6); // 璺宠繃 "data: " 鍓嶇紑
                        if (!"[DONE]".equals(data)) {
                            try {
                                ObjectMapper objectMapper = new ObjectMapper();
                                JsonNode rootNode = objectMapper.readTree(data);
                                JsonNode choicesNode = rootNode.path("choices");
                                if (choicesNode.isArray() && choicesNode.size() > 0) {
                                    JsonNode deltaNode = choicesNode.get(0).path("delta");
                                    JsonNode contentNode = deltaNode.path("content");
                                    if (contentNode != null && !contentNode.asText().isEmpty()) {
                                        String content = contentNode.asText();
                                        fullResponse.append(content);
                                        // 鍙戦€佹祦寮忔暟鎹埌鍓嶇锛岀‘淇濅娇鐢║TF-8缂栫爜
                                        String output = "data: " + objectMapper.writeValueAsString(content) + "\n\n";
                                        outputStream.write(output.getBytes(StandardCharsets.UTF_8));
                                        outputStream.flush();
                                    }
                                }
                            } catch (Exception e) {
                                // 蹇界暐瑙ｆ瀽閿欒锛岀户缁鐞嗕笅涓€琛?                                System.err.println("瑙ｆ瀽娴佸紡鍝嶅簲鍑洪敊: " + e.getMessage());
                            }
                        }
                    }
                }
            }
        } else {
            throw new RuntimeException("璋冪敤璞嗗寘API澶辫触: " + response.getStatusCode());
        }
        
        return fullResponse.toString();
    }
}
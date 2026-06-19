package com.long67.facetest.controller;

import com.long67.facetest.pojo.Conversation;
import com.long67.facetest.pojo.Message;
import com.long67.facetest.pojo.Result;
import com.long67.facetest.service.ConversationService;
import com.long67.facetest.utils.UserThreadLocal;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.mvc.method.annotation.StreamingResponseBody;

import java.io.OutputStream;
import java.util.List;

@RestController
@RequestMapping("/conversation")
public class ConversationController {
    
    @Autowired
    private ConversationService conversationService;
    
    /**
     * 鍒涘缓鏂扮殑瀵硅瘽
     * @param firstMessage 棣栨潯娑堟伅鍐呭
     * @return 瀵硅瘽瀵硅薄
     */
    @PostMapping("/create")
    public Result createConversation(@RequestParam String firstMessage) {
        try {
            // 浠嶶serThreadLocal鑾峰彇鐢ㄦ埛ID
            Integer userId = UserThreadLocal.getUserId();
            Conversation conversation = conversationService.createConversation(userId, firstMessage);
            return Result.success(conversation);
        } catch (Exception e) {
            return Result.error("鍒涘缓瀵硅瘽澶辫触: " + e.getMessage());
        }
    }
    
    /**
     * 鍙戦€佹秷鎭苟鑾峰彇AI鍥炲
     * @param conversationId 瀵硅瘽ID
     * @param content 娑堟伅鍐呭
     * @return AI鍥炲鍐呭
     */
    @PostMapping("/send")
    public Result sendMessage(@RequestParam Integer conversationId,
                            @RequestParam String content) {
        try {
            // 浠嶶serThreadLocal鑾峰彇鐢ㄦ埛ID
            Integer userId = UserThreadLocal.getUserId();
            String response = conversationService.sendMessage(conversationId, userId, content);
            return Result.success(response);
        } catch (Exception e) {
            return Result.error("鍙戦€佹秷鎭け璐? " + e.getMessage());
        }
    }
    
    /**
     * 鍙戦€佹秷鎭苟娴佸紡鑾峰彇AI鍥炲
     * @param conversationId 瀵硅瘽ID
     * @param content 娑堟伅鍐呭
     * @return 娴佸紡鍝嶅簲
     */
    @PostMapping(value = "/stream", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public ResponseEntity<StreamingResponseBody> streamMessage(@RequestParam Integer conversationId,
                                                               @RequestParam String content) {
        try {
            // 浠嶶serThreadLocal鑾峰彇鐢ㄦ埛ID
            Integer userId = UserThreadLocal.getUserId();
            
            StreamingResponseBody responseBody = outputStream -> {
                try {
                    conversationService.streamMessage(conversationId, userId, content, outputStream);
                } catch (Exception e) {
                    e.printStackTrace();
                }
            };
            
            return ResponseEntity.ok()
                    .contentType(MediaType.TEXT_EVENT_STREAM)
                    .header("Cache-Control", "no-cache")
                    .header("X-Accel-Buffering", "no")
                    .header("Connection", "keep-alive")
                    .header("Content-Type", "text/event-stream;charset=UTF-8")
                    .body(responseBody);
        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.internalServerError().build();
        }
    }
    
    /**
     * 鑾峰彇鐢ㄦ埛瀵硅瘽鍒楄〃
     * @return 瀵硅瘽鍒楄〃
     */
    @GetMapping("/list")
    public Result getConversations() {
        try {
            // 浠嶶serThreadLocal鑾峰彇鐢ㄦ埛ID
            Integer userId = UserThreadLocal.getUserId();
            List<Conversation> conversations = conversationService.getConversationsByUserId(userId);
            return Result.success(conversations);
        } catch (Exception e) {
            return Result.error("鑾峰彇瀵硅瘽鍒楄〃澶辫触: " + e.getMessage());
        }
    }
    
    /**
     * 鑾峰彇瀵硅瘽娑堟伅鍒楄〃
     * @param conversationId 瀵硅瘽ID
     * @return 娑堟伅鍒楄〃
     */
    @GetMapping("/messages")
    public Result getMessages(@RequestParam Integer conversationId) {
        try {
            List<Message> messages = conversationService.getMessagesByConversationId(conversationId);
            return Result.success(messages);
        } catch (Exception e) {
            return Result.error("鑾峰彇娑堟伅鍒楄〃澶辫触: " + e.getMessage());
        }
    }
}
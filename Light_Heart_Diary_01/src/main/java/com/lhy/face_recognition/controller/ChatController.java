package com.lhy.face_recognition.controller;

import com.lhy.face_recognition.entity.AiMessage;
import com.lhy.face_recognition.entity.AiSession;
import com.lhy.face_recognition.mapper.AiMessageMapper;
import com.lhy.face_recognition.mapper.AiSessionMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/chat")
@CrossOrigin(origins = {"http://localhost:5173", "http://localhost:5174"}, allowCredentials = "true")
public class ChatController {

    @Autowired
    private AiSessionMapper aiSessionMapper;

    @Autowired
    private AiMessageMapper aiMessageMapper;

    /**
     * 创建新会话
     */
    @PostMapping("/session/create")
    public ResponseEntity<Map<String, Object>> createSession(@RequestBody Map<String, Object> request) {
        try {
            Long userId = Long.valueOf(request.get("userId").toString());
            String sessionName = (String) request.getOrDefault("sessionName", "新对话");
            AiSession session = new AiSession();
            session.setUserId(userId);
            session.setSessionName(sessionName);
            aiSessionMapper.insertAiSession(session);
            return ResponseEntity.ok(Map.of("success", true, "session", session));
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(Map.of("success", false, "error", e.getMessage()));
        }
    }

    /**
     * 获取用户的会话列表
     */
    @GetMapping("/sessions/{userId}")
    public ResponseEntity<List<AiSession>> getSessions(@PathVariable Long userId) {
        return ResponseEntity.ok(aiSessionMapper.getSessionsByUserId(userId));
    }

    /**
     * 获取会话中的消息
     */
    @GetMapping("/messages/{sessionId}")
    public ResponseEntity<List<AiMessage>> getMessages(@PathVariable Long sessionId) {
        return ResponseEntity.ok(aiMessageMapper.getMessagesBySessionId(sessionId));
    }

    /**
     * 保存消息
     */
    @PostMapping("/message/save")
    public ResponseEntity<Map<String, Object>> saveMessage(@RequestBody AiMessage message) {
        try {
            aiMessageMapper.insertAiMessage(message);
            aiSessionMapper.updateSessionTime(message.getSessionId());
            return ResponseEntity.ok(Map.of("success", true, "message", message));
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(Map.of("success", false, "error", e.getMessage()));
        }
    }

    /**
     * 删除会话（连带消息）
     */
    @DeleteMapping("/session/{sessionId}")
    public ResponseEntity<Map<String, Object>> deleteSession(@PathVariable Long sessionId) {
        aiMessageMapper.deleteMessagesBySessionId(sessionId);
        aiSessionMapper.deleteSessionById(sessionId);
        return ResponseEntity.ok(Map.of("success", true));
    }
}
package com.lhy.face_recognition.controller;

import com.lhy.face_recognition.entity.DreamImage;
import com.lhy.face_recognition.service.DreamImageService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/dream-image")
@CrossOrigin(origins = {"http://localhost:5173", "http://localhost:5174"}, allowCredentials = "true")
public class DreamImageController {

    @Autowired
    private DreamImageService dreamImageService;

    /**
     * 保存绘梦记录
     */
    @PostMapping("/save")
    public ResponseEntity<Map<String, Object>> saveDreamImage(@RequestBody DreamImage dreamImage) {
        try {
            DreamImage saved = dreamImageService.saveDreamImage(dreamImage);
            return ResponseEntity.ok(Map.of("success", true, "data", saved));
        } catch (Exception e) {
            return ResponseEntity.badRequest().body(Map.of("success", false, "error", e.getMessage()));
        }
    }

    /**
     * 获取用户的绘梦历史
     */
    @GetMapping("/list/{userId}")
    public ResponseEntity<List<DreamImage>> getDreamImages(@PathVariable Long userId) {
        return ResponseEntity.ok(dreamImageService.getDreamImagesByUserId(userId));
    }

    /**
     * 删除绘梦记录
     */
    @DeleteMapping("/{id}/user/{userId}")
    public ResponseEntity<Map<String, Object>> deleteDreamImage(@PathVariable Long id, @PathVariable Long userId) {
        boolean deleted = dreamImageService.deleteDreamImage(id, userId);
        return ResponseEntity.ok(Map.of("success", deleted));
    }
}
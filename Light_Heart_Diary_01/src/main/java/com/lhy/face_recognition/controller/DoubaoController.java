package com.lhy.face_recognition.controller;

import com.lhy.face_recognition.service.DoubaoService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/doubao")
public class DoubaoController {

    @Autowired
    private DoubaoService doubaoService;

    /**
     * 调用豆包API获取回答
     * @param messages 聊天消息列表
     * @return 豆包API的回答
     */
    @PostMapping("/chat")
    public Map<String, Object> getChatAnswer(@RequestBody List<Map<String, Object>> messages) {
        try {
            String answer = doubaoService.getDoubaoAnswer(messages);
            return Map.of("success", true, "answer", answer);
        } catch (Exception e) {
            e.printStackTrace();
            return Map.of("success", false, "error", e.getMessage());
        }
    }
    
    /**
     * 调用豆包API生成图片
     * @param request 请求参数
     * @return 生成的图片数据
     */
    @PostMapping("/generate-image")
    public Map<String, Object> generateImage(@RequestBody Map<String, String> request) {
        try {
            String imageData = request.get("imageData");
            String mood = request.get("mood");
            
            if (imageData == null || mood == null) {
                return Map.of("success", false, "error", "缺少必要参数：imageData 或 mood");
            }
            
            String generatedImage = doubaoService.generateImage(imageData, mood);
            
            if (generatedImage != null && !generatedImage.isEmpty()) {
                return Map.of("success", true, "image", generatedImage);
            } else {
                // 服务层返回空，可能是内部逻辑判断失败
                return Map.of("success", false, "error", "图片生成服务返回为空，请检查日志");
            }
        } catch (Exception e) {
            e.printStackTrace(); // 控制台打印完整堆栈，方便调试
            
            // 提取更友好的错误信息
            String errorMsg = e.getMessage();
            if (errorMsg != null) {
                if (errorMsg.contains("InvalidParameter")) {
                    errorMsg = "API参数错误：模型不支持当前配置（如 output_format 或 size），请联系管理员修复。";
                } else if (errorMsg.contains("QuotaExceeded") || errorMsg.contains("额度")) {
                    errorMsg = "账户额度不足，请充值后重试。";
                } else if (errorMsg.contains("400")) {
                    errorMsg = "请求被服务器拒绝，请检查参数配置。";
                }
            }
            
            return Map.of("success", false, "error", errorMsg != null ? errorMsg : "未知错误，请稍后重试");
        }
    }
}

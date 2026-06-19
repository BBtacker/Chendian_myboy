package com.lhy.face_recognition.service;

import java.util.List;
import java.util.Map;

public interface DoubaoService {
    /**
     * 调用豆包API获取回答
     * @param messages 聊天消息列表
     * @return 豆包API的回答
     */
    String getDoubaoAnswer(List<Map<String, Object>> messages);
    
    /**
     * 调用豆包API生成图片
     * @param imageData 图片数据（base64）
     * @param mood 心情类型
     * @return 生成的图片数据（base64）
     */
    String generateImage(String imageData, String mood);
}

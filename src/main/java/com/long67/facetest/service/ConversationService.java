package com.long67.facetest.service;

import com.long67.facetest.pojo.Conversation;
import com.long67.facetest.pojo.Message;

import java.io.OutputStream;
import java.util.List;

public interface ConversationService {
    
    /**
     * 鍒涘缓鏂扮殑瀵硅瘽
     * @param userId 鐢ㄦ埛ID
     * @param firstMessage 棣栨潯娑堟伅鍐呭锛岀敤浜庣敓鎴愭爣棰?     * @return 瀵硅瘽瀵硅薄
     */
    Conversation createConversation(Integer userId, String firstMessage);
    
    /**
     * 鍙戦€佹秷鎭苟鑾峰彇AI鍥炲
     * @param conversationId 瀵硅瘽ID
     * @param userId 鐢ㄦ埛ID
     * @param content 鐢ㄦ埛娑堟伅鍐呭
     * @return AI鍥炲鍐呭
     * @throws Exception 澶勭悊寮傚父
     */
    String sendMessage(Integer conversationId, Integer userId, String content) throws Exception;
    
    /**
     * 娴佸紡鍙戦€佹秷鎭苟鑾峰彇AI鍥炲
     * @param conversationId 瀵硅瘽ID
     * @param userId 鐢ㄦ埛ID
     * @param content 鐢ㄦ埛娑堟伅鍐呭
     * @param outputStream 杈撳嚭娴?     * @throws Exception 澶勭悊寮傚父
     */
    void streamMessage(Integer conversationId, Integer userId, String content, OutputStream outputStream) throws Exception;
    
    /**
     * 鏍规嵁鐢ㄦ埛ID鏌ヨ瀵硅瘽鍒楄〃
     * @param userId 鐢ㄦ埛ID
     * @return 瀵硅瘽鍒楄〃
     */
    List<Conversation> getConversationsByUserId(Integer userId);
    
    /**
     * 鏍规嵁瀵硅瘽ID鏌ヨ娑堟伅鍒楄〃
     * @param conversationId 瀵硅瘽ID
     * @return 娑堟伅鍒楄〃
     */
    List<Message> getMessagesByConversationId(Integer conversationId);
}
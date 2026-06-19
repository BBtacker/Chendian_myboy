package com.long67.facetest.mapper;

import com.long67.facetest.pojo.Message;
import org.apache.ibatis.annotations.*;

import java.util.List;

@Mapper
public interface MessageMapper {
    
    /**
     * 鎻掑叆鏂扮殑娑堟伅璁板綍
     * @param message 娑堟伅瀵硅薄
     * @return 褰卞搷琛屾暟
     */
    @Insert("INSERT INTO message(conversation_id, sender, content, send_time) " +
            "VALUES(#{conversationId}, #{sender}, #{content}, #{sendTime})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertMessage(Message message);
    
    /**
     * 鏍规嵁瀵硅瘽ID鏌ヨ娑堟伅鍒楄〃
     * @param conversationId 瀵硅瘽ID
     * @return 娑堟伅鍒楄〃
     */
    @Select("SELECT * FROM message WHERE conversation_id = #{conversationId} ORDER BY send_time ASC")
    List<Message> selectMessagesByConversationId(Integer conversationId);
}
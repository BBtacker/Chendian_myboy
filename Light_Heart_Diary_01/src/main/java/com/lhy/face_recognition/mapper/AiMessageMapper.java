package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.AiMessage;
import org.apache.ibatis.annotations.*;

import java.util.List;

/**
 * AI消息Mapper接口
 */
@Mapper
public interface AiMessageMapper {

    @Insert("INSERT INTO ai_message (session_id, sender, content) VALUES (#{sessionId}, #{sender}, #{content})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertAiMessage(AiMessage aiMessage);

    @Select("SELECT * FROM ai_message WHERE session_id = #{sessionId} ORDER BY created_at ASC")
    List<AiMessage> getMessagesBySessionId(Long sessionId);

    @Delete("DELETE FROM ai_message WHERE session_id = #{sessionId}")
    int deleteMessagesBySessionId(Long sessionId);
}
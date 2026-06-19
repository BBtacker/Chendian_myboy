package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.AiSession;
import org.apache.ibatis.annotations.*;

import java.util.List;

/**
 * AI会话Mapper接口
 */
@Mapper
public interface AiSessionMapper {

    @Insert("INSERT INTO ai_session (user_id, session_name) VALUES (#{userId}, #{sessionName})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertAiSession(AiSession aiSession);

    @Select("SELECT * FROM ai_session WHERE user_id = #{userId} ORDER BY updated_at DESC")
    List<AiSession> getSessionsByUserId(Long userId);

    @Select("SELECT * FROM ai_session WHERE id = #{id}")
    AiSession getSessionById(Long id);

    @Update("UPDATE ai_session SET updated_at = NOW() WHERE id = #{id}")
    int updateSessionTime(Long id);

    @Delete("DELETE FROM ai_session WHERE id = #{id}")
    int deleteSessionById(Long id);
}
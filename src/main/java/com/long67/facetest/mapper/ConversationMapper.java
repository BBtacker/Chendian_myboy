package com.long67.facetest.mapper;

import com.long67.facetest.pojo.Conversation;
import org.apache.ibatis.annotations.*;

import java.util.List;

@Mapper
public interface ConversationMapper {
    
    /**
     * 鎻掑叆鏂扮殑瀵硅瘽璁板綍
     * @param conversation 瀵硅瘽瀵硅薄
     * @return 褰卞搷琛屾暟
     */
    @Insert("INSERT INTO conversation(user_id, title, start_time, last_update_time, status) " +
            "VALUES(#{userId}, #{title}, #{startTime}, #{lastUpdateTime}, #{status})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertConversation(Conversation conversation);
    
    /**
     * 鏍规嵁ID鏇存柊瀵硅瘽
     * @param conversation 瀵硅瘽瀵硅薄
     * @return 褰卞搷琛屾暟
     */
    @Update("UPDATE conversation SET title = #{title}, last_update_time = #{lastUpdateTime}, status = #{status} " +
            "WHERE id = #{id}")
    int updateConversation(Conversation conversation);
    
    /**
     * 鏍规嵁鐢ㄦ埛ID鏌ヨ瀵硅瘽鍒楄〃
     * @param userId 鐢ㄦ埛ID
     * @return 瀵硅瘽鍒楄〃
     */
    @Select("SELECT * FROM conversation WHERE user_id = #{userId} ORDER BY last_update_time DESC")
    List<Conversation> selectConversationsByUserId(Integer userId);
    
    /**
     * 鏍规嵁ID鏌ヨ瀵硅瘽
     * @param id 瀵硅瘽ID
     * @return 瀵硅瘽瀵硅薄
     */
    @Select("SELECT * FROM conversation WHERE id = #{id}")
    Conversation selectConversationById(Integer id);
}
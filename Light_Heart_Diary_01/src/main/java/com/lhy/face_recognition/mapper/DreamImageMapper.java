package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.DreamImage;
import org.apache.ibatis.annotations.*;

import java.util.List;

/**
 * 绘梦记录Mapper接口
 */
@Mapper
public interface DreamImageMapper {

    @Insert("INSERT INTO dream_image (user_id, prompt, mood, image_url, original_image) " +
            "VALUES (#{userId}, #{prompt}, #{mood}, #{imageUrl}, #{originalImage})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertDreamImage(DreamImage dreamImage);

    @Select("SELECT * FROM dream_image WHERE user_id = #{userId} ORDER BY created_at DESC")
    List<DreamImage> getDreamImagesByUserId(Long userId);

    @Select("SELECT * FROM dream_image WHERE id = #{id}")
    DreamImage getDreamImageById(Long id);

    @Delete("DELETE FROM dream_image WHERE id = #{id} AND user_id = #{userId}")
    int deleteDreamImage(@Param("id") Long id, @Param("userId") Long userId);
}
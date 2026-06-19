package com.lhy.face_recognition.mapper;

import com.lhy.face_recognition.entity.DiaryPhoto;
import org.apache.ibatis.annotations.*;

import java.util.List;

/**
 * 日记照片Mapper接口
 */
@Mapper
public interface DiaryPhotoMapper {
    /**
     * 插入日记照片
     * @param diaryPhoto 日记照片实体
     * @return 影响的行数
     */
    @Insert("INSERT INTO diary_photo (diary_id, photo_url, created_at) VALUES (#{diaryId}, #{photoUrl}, #{createdAt})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertDiaryPhoto(DiaryPhoto diaryPhoto);
    
    /**
     * 根据日记ID获取照片列表
     * @param diaryId 日记ID
     * @return 照片列表
     */
    @Select("SELECT * FROM diary_photo WHERE diary_id = #{diaryId} ORDER BY id")
    List<DiaryPhoto> getPhotosByDiaryId(Long diaryId);
    
    /**
     * 根据ID删除照片
     * @param id 照片ID
     * @return 影响的行数
     */
    @Delete("DELETE FROM diary_photo WHERE id = #{id}")
    int deletePhotoById(Long id);
    
    /**
     * 根据日记ID删除所有照片
     * @param diaryId 日记ID
     * @return 影响的行数
     */
    @Delete("DELETE FROM diary_photo WHERE diary_id = #{diaryId}")
    int deletePhotosByDiaryId(Long diaryId);
}
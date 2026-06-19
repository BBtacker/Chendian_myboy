package com.long67.facetest.mapper;

import com.long67.facetest.pojo.User;
import org.apache.ibatis.annotations.*;

import java.util.List;

@Mapper
public interface UserMapper {
    /**
     * 鏌ヨ鎵€鏈夌敤鎴?
     */
    @Select("SELECT * FROM user")
    List<User> list();

    /**
     * 鏍规嵁ID鏌ヨ鐢ㄦ埛
     */
    @Select("SELECT * FROM user WHERE id = #{id}")
    User getById(Integer id);

    /**
     * 鏍规嵁鐢ㄦ埛鍚嶆煡璇㈢敤鎴?
     */
    @Select("SELECT * FROM user WHERE username = #{username}")
    User getByUsername(String username);

    /**
     * 娣诲姞鐢ㄦ埛
     */
    @Insert("INSERT INTO user(username, password, create_time, avatar, name, email, phone, gender, birthday, address) VALUES(#{username}, #{password}, #{createTime}, #{avatar}, #{name}, #{email}, #{phone}, #{gender}, #{birthday}, #{address})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    void add(User user);

    /**
     * 鏇存柊鐢ㄦ埛淇℃伅锛堜笉鏇存柊瀵嗙爜锛?
     */
    @Update("UPDATE user SET username = #{username}, avatar = #{avatar}, name = #{name}, email = #{email}, phone = #{phone}, gender = #{gender}, birthday = #{birthday}, address = #{address} WHERE id = #{id}")
    void update(User user);
    
    /**
     * 鏇存柊鐢ㄦ埛瀵嗙爜
     */
    @Update("UPDATE user SET password = #{password} WHERE id = #{id}")
    void updatePassword(@Param("id") Integer id, @Param("password") String password);

    /**
     * 鏍规嵁ID鍒犻櫎鐢ㄦ埛
     */
    @Delete("DELETE FROM user WHERE id = #{id}")
    void deleteById(Integer id);
    
    /**
     * 鏉′欢鏌ヨ鐢ㄦ埛
     */
    @Select("SELECT * FROM user WHERE username LIKE CONCAT('%', #{username}, '%') OR #{username} IS NULL")
    List<User> getUsersByCondition(@Param("username") String username);
}
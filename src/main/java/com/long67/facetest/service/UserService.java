package com.long67.facetest.service;

import com.long67.facetest.pojo.User;
import com.github.pagehelper.PageInfo;

import java.util.List;

public interface UserService {
    /**
     * 鏌ヨ鎵€鏈夌敤鎴?
     */
    List<User> list();

    /**
     * 鏍规嵁ID鏌ヨ鐢ㄦ埛
     */
    User getById(Integer id);

    /**
     * 鏍规嵁鐢ㄦ埛鍚嶆煡璇㈢敤鎴?
     */
    User getByUsername(String username);

    /**
     * 娣诲姞鐢ㄦ埛
     */
    void add(User user);

    /**
     * 鏇存柊鐢ㄦ埛淇℃伅
     */
    void update(User user);

    /**
     * 鏍规嵁ID鍒犻櫎鐢ㄦ埛
     */
    void deleteById(Integer id);
    
    /**
     * 鍒嗛〉鏉′欢鏌ヨ鐢ㄦ埛
     */
    PageInfo<User> getUsersByPage(Integer page, Integer pageSize, String username);
}
package com.long67.facetest.service.Impl;

import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.long67.facetest.mapper.UserMapper;
import com.long67.facetest.pojo.User;
import com.long67.facetest.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.util.List;

@Service
public class UserServiceImpl implements UserService {
    
    @Autowired
    private UserMapper userMapper;
    
    /**
     * 鏌ヨ鎵€鏈夌敤鎴?
     */
    @Override
    public List<User> list() {
        return userMapper.list();
    }
    
    /**
     * 鏍规嵁ID鏌ヨ鐢ㄦ埛
     */
    @Override
    public User getById(Integer id) {
        return userMapper.getById(id);
    }
    
    /**
     * 鏍规嵁鐢ㄦ埛鍚嶆煡璇㈢敤鎴?
     */
    @Override
    public User getByUsername(String username) {
        return userMapper.getByUsername(username);
    }
    
    /**
     * 娣诲姞鐢ㄦ埛
     */
    @Override
    public void add(User user) {
        userMapper.add(user);
    }
    
    /**
     * 鏇存柊鐢ㄦ埛淇℃伅
     */
    @Override
    public void update(User user) {
        // 濡傛灉瀵嗙爜涓嶄负绌猴紝鍒欏崟鐙洿鏂板瘑鐮?
        if (StringUtils.hasText(user.getPassword())) {
            userMapper.updatePassword(user.getId(), user.getPassword());
        }
        // 鏇存柊鍏朵粬鐢ㄦ埛淇℃伅
        userMapper.update(user);
    }
    
    /**
     * 鏍规嵁ID鍒犻櫎鐢ㄦ埛
     */
    @Override
    public void deleteById(Integer id) {
        userMapper.deleteById(id);
    }
    
    /**
     * 鍒嗛〉鏉′欢鏌ヨ鐢ㄦ埛
     */
    @Override
    public PageInfo<User> getUsersByPage(Integer page, Integer pageSize, String username) {
        PageHelper.startPage(page, pageSize);
        List<User> users = userMapper.getUsersByCondition(username);
        return new PageInfo<>(users);
    }
}
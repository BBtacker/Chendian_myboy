package com.lhy.face_recognition.service;

import com.lhy.face_recognition.entity.User;

/**
 * 用户服务接口
 */
public interface UserService {
    /**
     * 用户登录
     * @param username 用户名
     * @param password 密码
     * @return 登录成功返回用户信息，失败返回null
     */
    User login(String username, String password);

    /**
     * 用户注册
     * @param username 用户名
     * @param password 密码
     * @return 注册成功返回用户信息，失败返回null
     */
    User register(String username, String password);

    /**
     * 根据ID获取用户信息
     * @param id 用户ID
     * @return 用户信息
     */
    User getUserById(Long id);

    /**
     * 更新用户信息
     * @param user 用户信息
     * @return 更新后的用户信息
     */
    User updateUser(User user);
}

package com.lhy.face_recognition.service.impl;

import com.lhy.face_recognition.entity.User;
import com.lhy.face_recognition.mapper.UserMapper;
import com.lhy.face_recognition.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

/**
 * 用户服务实现类
 */
@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private UserMapper userMapper;

    private final BCryptPasswordEncoder passwordEncoder = new BCryptPasswordEncoder();

    @Override
    public User login(String username, String password) {
        // 根据用户名查询用户
        User user = userMapper.findByUsername(username);
        if (user == null) {
            return null; // 用户不存在
        }

        // 验证密码
        if (passwordEncoder.matches(password, user.getPassword())) {
            // 登录成功，返回用户信息（不包含密码）
            user.setPassword(null);
            return user;
        }
        return null; // 密码错误
    }

    @Override
    public User register(String username, String password) {
        // 检查用户名是否已存在
        if (userMapper.findByUsername(username) != null) {
            return null; // 用户名已存在
        }

        // 加密密码
        String encodedPassword = passwordEncoder.encode(password);

        // 创建新用户
        User user = new User();
        user.setUsername(username);
        user.setPassword(encodedPassword);

        // 保存用户到数据库
        if (userMapper.insert(user) > 0) {
            // 注册成功，返回用户信息（不包含密码）
            user.setPassword(null);
            return user;
        }
        return null; // 注册失败
    }

    @Override
    public User getUserById(Long id) {
        User user = userMapper.findById(id);
        if (user != null) {
            user.setPassword(null); // 不返回密码
        }
        return user;
    }

    @Override
    public User updateUser(User user) {
        if (userMapper.update(user) > 0) {
            return getUserById(user.getId());
        }
        return null;
    }
}

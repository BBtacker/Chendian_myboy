package com.lhy.face_recognition.service;

import org.springframework.web.multipart.MultipartFile;

/**
 * 阿里云OSS服务接口
 */
public interface OssService {

    /**
     * 上传图片到OSS
     * @param file 图片文件
     * @return 图片的访问URL
     */
    String uploadImage(MultipartFile file);

    /**
     * 删除OSS上的图片
     * @param imageUrl 图片URL
     * @return 删除是否成功
     */
    boolean deleteImage(String imageUrl);
}

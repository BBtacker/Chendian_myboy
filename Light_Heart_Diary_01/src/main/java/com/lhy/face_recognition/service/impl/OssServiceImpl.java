package com.lhy.face_recognition.service.impl;

import com.aliyun.oss.OSS;
import com.aliyun.oss.model.ObjectMetadata;
import com.lhy.face_recognition.config.OssConfig;
import com.lhy.face_recognition.service.OssService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.io.InputStream;
import java.util.UUID;

/**
 * 阿里云OSS服务实现类
 */
@Service
public class OssServiceImpl implements OssService {

    @Autowired
    private OSS ossClient;

    @Autowired
    private OssConfig ossConfig;

    /**
     * 上传图片到OSS
     * @param file 图片文件
     * @return 图片的访问URL
     */
    @Override
    public String uploadImage(MultipartFile file) {
        if (file.isEmpty()) {
            throw new RuntimeException("上传文件不能为空");
        }

        try {
            // 获取文件名
            String originalFilename = file.getOriginalFilename();
            // 获取文件后缀
            String fileSuffix = originalFilename.substring(originalFilename.lastIndexOf("."));
            // 生成唯一文件名
            String fileName = "diary/" + UUID.randomUUID() + fileSuffix;

            // 获取文件输入流
            InputStream inputStream = file.getInputStream();

            // 设置文件元数据
            ObjectMetadata metadata = new ObjectMetadata();
            metadata.setContentType(file.getContentType());
            metadata.setContentLength(file.getSize());

            // 上传文件到OSS
            ossClient.putObject(ossConfig.getBucketName(), fileName, inputStream, metadata);

            // 关闭输入流
            inputStream.close();

            // 构造访问URL
            String url = "https://" + ossConfig.getBucketName() + "." + ossConfig.getEndpoint() + "/" + fileName;

            // 注意：这里返回的是模拟URL，实际使用时需要返回真实的OSS URL
            return url;
        } catch (IOException e) {
            throw new RuntimeException("上传文件失败：" + e.getMessage(), e);
        }
    }

    /**
     * 删除OSS上的图片
     * @param imageUrl 图片URL
     * @return 删除是否成功
     */
    @Override
    public boolean deleteImage(String imageUrl) {
        if (imageUrl == null || imageUrl.isEmpty()) {
            return false;
        }

        try {
            // 从URL中提取文件名
            String fileName = imageUrl.substring(imageUrl.lastIndexOf("/") + 1);
            // 删除OSS上的文件
            ossClient.deleteObject(ossConfig.getBucketName(), "diary/" + fileName);
            return true;
        } catch (Exception e) {
            throw new RuntimeException("删除文件失败：" + e.getMessage(), e);
        }
    }
}

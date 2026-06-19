package com.lhy.face_recognition.controller;

import com.lhy.face_recognition.service.OssService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.util.HashMap;
import java.util.Map;

/**
 * 阿里云OSS控制器
 */
@RestController
@RequestMapping("/api/oss")
public class OssController {

    @Autowired
    private OssService ossService;

    /**
     * 上传图片
     * @param file 图片文件
     * @return 上传结果
     */
    @PostMapping("/upload")
    public ResponseEntity<Map<String, Object>> uploadImage(@RequestParam("file") MultipartFile file) {
        try {
            // 调用OSS服务上传图片
            String imageUrl = ossService.uploadImage(file);

            // 构建响应结果
            Map<String, Object> result = new HashMap<>();
            result.put("success", true);
            result.put("message", "图片上传成功");
            result.put("data", imageUrl);

            return new ResponseEntity<>(result, HttpStatus.OK);
        } catch (Exception e) {
            Map<String, Object> error = new HashMap<>();
            error.put("success", false);
            error.put("message", "图片上传失败：" + e.getMessage());
            return new ResponseEntity<>(error, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    /**
     * 删除图片
     * @param imageUrl 图片URL
     * @return 删除结果
     */
    @DeleteMapping("/delete")
    public ResponseEntity<Map<String, Object>> deleteImage(@RequestParam("imageUrl") String imageUrl) {
        try {
            // 调用OSS服务删除图片
            boolean success = ossService.deleteImage(imageUrl);

            // 构建响应结果
            Map<String, Object> result = new HashMap<>();
            result.put("success", success);
            result.put("message", success ? "图片删除成功" : "图片删除失败");

            return new ResponseEntity<>(result, HttpStatus.OK);
        } catch (Exception e) {
            Map<String, Object> error = new HashMap<>();
            error.put("success", false);
            error.put("message", "图片删除失败：" + e.getMessage());
            return new ResponseEntity<>(error, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}

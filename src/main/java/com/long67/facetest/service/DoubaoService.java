package com.long67.facetest.service;

import com.long67.facetest.pojo.testResult;
import org.springframework.web.multipart.MultipartFile;

public interface DoubaoService {
    /**
     * 鍒嗘瀽闈㈤儴鍥剧墖鏄惁涓鸿吅浣撻潰瀹?
     * @param image 鍥剧墖鏂囦欢
     * @param userId 鐢ㄦ埛ID
     * @return 鍒嗘瀽缁撴灉
     * @throws Exception 澶勭悊寮傚父
     */
    testResult analyzeFace(MultipartFile image, Integer userId) throws Exception;
}
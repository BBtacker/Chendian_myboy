package com.lhy.face_recognition.service;

import com.lhy.face_recognition.entity.DreamImage;
import java.util.List;

public interface DreamImageService {
    DreamImage saveDreamImage(DreamImage dreamImage);
    List<DreamImage> getDreamImagesByUserId(Long userId);
    DreamImage getDreamImageById(Long id);
    boolean deleteDreamImage(Long id, Long userId);
}
package com.lhy.face_recognition.service.impl;

import com.lhy.face_recognition.entity.DreamImage;
import com.lhy.face_recognition.mapper.DreamImageMapper;
import com.lhy.face_recognition.service.DreamImageService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class DreamImageServiceImpl implements DreamImageService {

    @Autowired
    private DreamImageMapper dreamImageMapper;

    @Override
    public DreamImage saveDreamImage(DreamImage dreamImage) {
        dreamImageMapper.insertDreamImage(dreamImage);
        return dreamImage;
    }

    @Override
    public List<DreamImage> getDreamImagesByUserId(Long userId) {
        return dreamImageMapper.getDreamImagesByUserId(userId);
    }

    @Override
    public DreamImage getDreamImageById(Long id) {
        return dreamImageMapper.getDreamImageById(id);
    }

    @Override
    public boolean deleteDreamImage(Long id, Long userId) {
        return dreamImageMapper.deleteDreamImage(id, userId) > 0;
    }
}
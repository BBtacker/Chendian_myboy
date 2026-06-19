package com.long67.facetest.pojo;

import lombok.Data;

import java.time.LocalDateTime;

@Data
public class Message {
    private Integer id;
    private Integer conversationId;
    private Integer sender; // 0=鐢ㄦ埛锛?=璞嗗寘鍔╂墜
    private String content;
    private LocalDateTime sendTime;
}
package com.long67.facetest.pojo;

import lombok.Data;
import java.time.LocalDateTime;
import java.time.LocalDate;

@Data
public class User {
    private Integer id;
    private String username;
    private String password;
    private LocalDateTime createTime;
    private String avatar;
    private String name;
    private String email;
    private String phone;
    private Integer gender; // 0-鏈煡锛?-鐢凤紝2-濂?
    private LocalDate birthday;
    private String address;
}
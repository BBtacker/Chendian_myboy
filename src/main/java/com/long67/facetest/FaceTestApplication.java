package com.long67.facetest;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@MapperScan("com.long67.facetest.mapper")
public class FaceTestApplication {

    public static void main(String[] args) {
        SpringApplication.run(FaceTestApplication.class, args);
    }

}
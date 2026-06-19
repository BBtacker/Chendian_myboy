package com.long67.facetest.pojo;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Data;

import java.time.LocalDateTime;

@Data
public class testResult {
    private Integer id;
    private Integer userId;
    private String imagePath;
    private Boolean isGlandFace;
    private String level;
    private Double confidence;
    private String visualizationDescription;
    
    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime createTime;
    
    // 鍏煎鏁版嵁搴撳瓧娈靛悕
    private LocalDateTime testTime;
    
    // 鐢ㄤ簬搴忓垪鍖栨椂杩斿洖createTime瀛楁
    public LocalDateTime getCreateTime() {
        return testTime;
    }
}
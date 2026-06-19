package com.long67.facetest.pojo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.util.List;

/**
 * 灏佽鍒嗛〉鏌ヨ缁撴灉
 */
@Data
@AllArgsConstructor
@NoArgsConstructor
public class PageResult implements Serializable {

    private long total; //鎬昏褰曟暟

    private List records; //褰撳墠椤垫暟鎹泦鍚?

}

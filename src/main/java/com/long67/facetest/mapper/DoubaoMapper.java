package com.long67.facetest.mapper;

import com.long67.facetest.pojo.testResult;
import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Options;

@Mapper
public interface DoubaoMapper {
    
    /**
     * 鎻掑叆娴嬭瘯缁撴灉
     * @param testResult 娴嬭瘯缁撴灉瀵硅薄
     * @return 褰卞搷琛屾暟
     */
    @Insert("INSERT INTO test_result(user_id, image, test_time, is_gland_face, level, probability, visualization_description) " +
            "VALUES(#{userId}, #{imagePath}, #{testTime}, #{isGlandFace}, #{level}, #{confidence}, #{visualizationDescription})")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    int insertTestResult(testResult testResult);
}
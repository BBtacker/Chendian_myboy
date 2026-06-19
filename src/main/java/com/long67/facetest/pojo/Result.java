package com.long67.facetest.pojo;

import lombok.Data;

import java.io.Serializable;

/**
 * 鍚庣缁熶竴杩斿洖缁撴灉
 * @param <T>
 */
@Data
public class Result<T> implements Serializable {

    private Integer code; //缂栫爜锛?鎴愬姛锛?鍜屽叾瀹冩暟瀛椾负澶辫触
    private String msg; //閿欒淇℃伅
    private T data; //鏁版嵁

    public static <T> Result<T> success() {
        Result<T> result = new Result<T>();
        result.code = 1;
        return result;
    }

    public static <T> Result<T> success(T object) {
        Result<T> result = new Result<T>();
        result.data = object;
        result.code = 1;
        return result;
    }

    public static <T> Result<T> error(String msg) {
        Result<T> result = new Result<T>();
        result.msg = msg;
        result.code = 0;
        return result;
    }
}
"""
特征提取服务入口 (FastAPI)
==========================
接口:
  POST /extract     multipart 上传图片 -> {"feature_vector": [...1536], "dimension": 1536, "description": {...}, "face_detected": bool, "measurements": {...}}
  GET  /health      健康检查

启动: uvicorn app.main:app --host 0.0.0.0 --port 8085
"""

import io
import time

from fastapi import FastAPI, File, UploadFile
from fastapi.responses import JSONResponse

from app import extractor

app = FastAPI(title="Feature Extractor", version="1.0.0")

# 启动时预热加载模型（避免首个请求慢）
@app.on_event("startup")
def warmup():
    try:
        extractor.load_model()
    except Exception as e:
        print(f"[warn] 模型预热失败（服务仍可启动，请求时会报错）: {e}", flush=True)
    try:
        from app import facemesh
        facemesh.warmup()
    except Exception as e:
        print(f"[warn] FaceMesh 预热失败: {e}", flush=True)


@app.get("/health")
def health():
    return {
        "status": "ok",
        "model_loaded": extractor._session is not None,
        "model_path": extractor.MODEL_PATH,
        "feature_dim": extractor.FEATURE_DIM,
    }


@app.post("/extract")
async def extract(image: UploadFile = File(...)):
    """提取图片特征向量"""
    start = time.time()

    # 校验文件类型（content-type 或扩展名任一匹配即可，避免 Go 端 multipart 未带 content-type 被拒）
    _allowed_ext = (".jpg", ".jpeg", ".png", ".webp", ".gif", ".bmp")
    _fname = (image.filename or "").lower()
    _ct = (image.content_type or "").lower()
    if not (_ct.startswith("image/") or _fname.endswith(_allowed_ext)):
        return JSONResponse(status_code=400, content={"code": 0, "msg": "不支持的图片类型: " + str(image.content_type)})

    # 读取并解码
    data = await image.read()
    if len(data) == 0:
        return JSONResponse(status_code=400, content={"code": 0, "msg": "图片内容为空"})

    try:
        from PIL import Image
        img = Image.open(io.BytesIO(data))
        img.load()
        img = img.convert("RGB")
    except Exception as e:
        return JSONResponse(status_code=400, content={"code": 0, "msg": f"图片解码失败: {e}"})

    # 特征提取
    try:
        feature = extractor.extract_feature(img)
    except Exception as e:
        return JSONResponse(status_code=500, content={"code": 0, "msg": f"特征提取失败: {e}"})

    # 特征描述（基于 B3 向量的确定性统计，用于 RAG 相似检索）
    description = extractor.generate_description(feature)

    # 真实面部几何测量（MediaPipe Face Mesh，可选；不可用时不阻断主流程）
    # 关键区分："人脸门控不可用" vs "确实未检测到人脸"
    #   - facemesh 正常跑且报无脸 -> face_detected=False（真实剔除，符合需求）
    #   - facemesh 不可用（如未装 mediapipe）-> 无法判断，按 lenient 视为有人脸，避免误剔除全部照片
    face_detected = True
    measurements: dict = {}
    try:
        from app import facemesh
        mesh = facemesh.extract_measurements(img)
        face_detected = mesh.get("face_detected", True)
        measurements = mesh.get("measurements", {})
    except Exception as e:
        print(f"[warn] FaceMesh 不可用，按有人脸处理（安装 mediapipe 后启用真实人脸门控）: {e}", flush=True)

    return {
        "code": 1,
        "msg": "success",
        "feature_vector": feature.tolist(),
        "dimension": int(extractor.FEATURE_DIM),
        "description": description,
        "face_detected": face_detected,
        "measurements": measurements,
        "elapsed_ms": int((time.time() - start) * 1000),
    }

"""
MediaPipe Face Mesh 真实面部几何测量
===================================
从人脸 478 个 landmarks 计算腺样体面容相关的客观几何测量值，
作为 DeepSeek 诊断 prompt 的"客观辅助证据"，替代原先基于向量统计的占位描述。

计算项（均为归一化/比例量，与拍摄距离无关）：
  - mouth_aperture_ratio       张口度/口宽比（越大越倾向张口呼吸）
  - lower_face_ratio           面下1/3比例（腺样体面容常增大）
  - face_width_height_ratio    面宽/面高比（偏小=长面型）
  - nasolabial_angle_deg       鼻唇角（度）
  - facial_convexity_angle_deg 面凸角（度，下颌后缩 proxy）

依赖: pip install mediapipe
模型: MediaPipe Face Mesh 模型在首次调用时自动联网下载（无需手动准备）。
"""

import math

# MediaPipe Face Mesh 规范 landmark 索引
LANDMARKS = {
    "forehead": 10,      # 发际中点（面高上界）
    "nasion": 168,       # 鼻根（两眼之间）
    "nose_tip": 1,       # 鼻尖
    "subnasale": 2,      # 鼻下点（鼻柱基底）
    "upper_lip": 0,      # 上唇中点
    "lower_lip": 17,     # 下唇中点
    "mouth_left": 61,    # 口裂左角
    "mouth_right": 291,  # 口裂右角
    "chin": 152,         # 颏下点（面下界）
    "cheek_left": 234,   # 左颧弓最外
    "cheek_right": 454,  # 右颧弓最外
}


def _xy(lm, idx, w, h):
    p = lm[idx]
    # MediaPipe landmark 坐标为相对图像的归一化值，乘回像素尺寸得到真实距离
    return p.x * w, p.y * h


def _dist(a, b):
    return math.hypot(a[0] - b[0], a[1] - b[1])


def _angle(p, a, b):
    """以 p 为顶点，向量 p->a 与 p->b 的夹角（度）"""
    v1 = (a[0] - p[0], a[1] - p[1])
    v2 = (b[0] - p[0], b[1] - p[1])
    dot = v1[0] * v2[0] + v1[1] * v2[1]
    n1 = math.hypot(*v1)
    n2 = math.hypot(*v2)
    if n1 == 0 or n2 == 0:
        return 0.0
    cosv = max(-1.0, min(1.0, dot / (n1 * n2)))
    return math.degrees(math.acos(cosv))


_FACE_MESH = None


def _get_face_mesh():
    global _FACE_MESH
    if _FACE_MESH is None:
        import mediapipe as mp
        mp_face_mesh = mp.solutions.face_mesh
        _FACE_MESH = mp_face_mesh.FaceMesh(
            static_image_mode=True,
            max_num_faces=1,
            refine_landmarks=True,
            min_detection_confidence=0.5,
        )
    return _FACE_MESH


def warmup():
    """应用启动时预加载 Face Mesh 模型，避免首次请求才联网下载导致特征服务超时。"""
    try:
        _get_face_mesh()
        print("[warmup] FaceMesh 模型已就绪", flush=True)
        return True
    except Exception as e:
        print(f"[warn] FaceMesh 预热失败（首次 /extract 会触发下载）: {e}", flush=True)
        return False


def extract_measurements(image):
    """从 PIL.Image 计算面部几何测量值。

    返回 dict: {"face_detected": bool, "measurements": {...}, "landmark_count": int}
    未检测到人脸时 face_detected=False, measurements={}。
    """
    import numpy as np
    from PIL import Image

    rgb = image.convert("RGB")
    arr = np.asarray(rgb)
    h, w = arr.shape[:2]

    fm = _get_face_mesh()
    results = fm.process(arr)
    if not results.multi_face_landmarks:
        return {"face_detected": False, "measurements": {}, "landmark_count": 0}

    lm = results.multi_face_landmarks[0].landmark
    n = len(lm)
    pts = {k: _xy(lm, idx, w, h) for k, idx in LANDMARKS.items()}

    mouth_w = _dist(pts["mouth_left"], pts["mouth_right"])
    mouth_gap = _dist(pts["upper_lip"], pts["lower_lip"])
    face_h = _dist(pts["forehead"], pts["chin"])
    lower_face = _dist(pts["subnasale"], pts["chin"])
    face_w = _dist(pts["cheek_left"], pts["cheek_right"])

    measurements = {
        "mouth_aperture_ratio": round(mouth_gap / mouth_w, 4) if mouth_w > 0 else 0.0,
        "lower_face_ratio": round(lower_face / face_h, 4) if face_h > 0 else 0.0,
        "face_width_height_ratio": round(face_w / face_h, 4) if face_h > 0 else 0.0,
        "nasolabial_angle_deg": round(_angle(pts["subnasale"], pts["nose_tip"], pts["upper_lip"]), 2),
        "facial_convexity_angle_deg": round(_angle(pts["subnasale"], pts["nasion"], pts["chin"]), 2),
    }
    return {"face_detected": True, "measurements": measurements, "landmark_count": n}

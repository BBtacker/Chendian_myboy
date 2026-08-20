"""
EfficientNet-B3 特征提取器
=========================
基于 ONNX Runtime 的 EfficientNet-B3 特征提取：
1. 加载 efficientnet_b3.onnx（去掉分类头的特征层，输出 1536 维特征向量）
2. 图片预处理：resize 300x300 -> normalize（ImageNet 标准）
3. 推理输出 1536 维 float32 特征向量（L2 归一化）
4. 可选生成特征描述（stats 模式：基于向量的确定性统计，占位实现）

模型文件: models/efficientnet_b3.onnx
获取方式: 见 README.md（torchvision 导出 或 下载现成 ONNX）
"""

import json
import os
import time

import numpy as np
from PIL import Image

# ---- 配置（优先环境变量，其次 config.env）----
MODEL_PATH = os.environ.get("MODEL_PATH", "models/efficientnet_b3.onnx")
INPUT_SIZE = int(os.environ.get("INPUT_SIZE", "300"))
FEATURE_DIM = int(os.environ.get("FEATURE_DIM", "1536"))
DESCRIPTION_MODE = os.environ.get("DESCRIPTION_MODE", "stats")

# ImageNet 归一化参数（torchvision 标准）
MEAN = np.array([0.485, 0.456, 0.406], dtype=np.float32)
STD = np.array([0.229, 0.224, 0.225], dtype=np.float32)

_session = None
_input_name = None


def load_model():
    """加载 ONNX 模型（惰性加载，首次请求时初始化）"""
    global _session, _input_name
    if _session is not None:
        return _session

    if not os.path.exists(MODEL_PATH):
        raise FileNotFoundError(
            f"模型文件不存在: {MODEL_PATH}。请下载 efficientnet_b3.onnx 放入 models/ 目录，"
            "获取方式见 feature_extractor/README.md"
        )

    import onnxruntime as ort

    # CPU 推理（可按需启用 GPU: providers=['CUDAExecutionProvider', 'CPUExecutionProvider']）
    _session = ort.InferenceSession(
        MODEL_PATH, providers=["CPUExecutionProvider"]
    )
    _input_name = _session.get_inputs()[0].name

    # 校验输出维度
    output_shape = _session.get_outputs()[0].shape
    if output_shape and len(output_shape) == 2 and output_shape[1] != FEATURE_DIM:
        print(f"[warn] 模型输出维度 {output_shape[1]} 与配置 {FEATURE_DIM} 不一致，"
              "请检查模型是否已去掉分类头", flush=True)

    print(f"[OK] 模型加载成功: {MODEL_PATH} (输入: {_input_name})", flush=True)
    return _session


def preprocess(image: Image.Image) -> np.ndarray:
    """图片预处理：resize -> RGB -> CHW -> normalize"""
    # 保持宽高比缩放后中心裁剪（torchvision 默认行为）
    w, h = image.size
    shorter = min(w, h)
    image = image.crop(((w - shorter) // 2, (h - shorter) // 2,
                        (w + shorter) // 2, (h + shorter) // 2))
    image = image.resize((INPUT_SIZE, INPUT_SIZE), Image.BILINEAR)

    arr = np.asarray(image, dtype=np.float32)
    if arr.ndim == 2:  # 灰度图转 RGB
        arr = np.stack([arr] * 3, axis=-1)
    elif arr.shape[2] == 4:  # RGBA 去 alpha
        arr = arr[:, :, :3]

    # HWC -> CHW 并归一化
    arr = arr.transpose(2, 0, 1)
    arr = (arr / 255.0 - MEAN[:, None, None]) / STD[:, None, None]
    return arr.astype(np.float32)


def extract_feature(image: Image.Image) -> np.ndarray:
    """提取 1536 维特征向量（L2 归一化）"""
    session = load_model()
    tensor = preprocess(image)
    # 添加 batch 维度: (1, 3, 300, 300)
    tensor = tensor[np.newaxis, ...]

    outputs = session.run(None, {_input_name: tensor})
    feature = np.asarray(outputs[0]).reshape(-1)

    # 若模型带分类头（1000 类），截取特征部分不可行——需要去掉分类头的模型。
    # 这里校验长度，异常时提示用户检查模型。
    if feature.shape[0] != FEATURE_DIM:
        raise ValueError(
            f"模型输出维度 {feature.shape[0]} != {FEATURE_DIM}，"
            "请使用去掉分类头(Classifier)的 EfficientNet-B3 ONNX 模型。"
            "导出方式: torchvision.models.efficientnet_b3(weights=IMAGENET1K_V1), "
            "删除 features 后的 classifier，只导出 backbone"
        )

    # L2 归一化（保证向量比较有意义的尺度）
    norm = np.linalg.norm(feature)
    if norm > 0:
        feature = feature / norm
    return feature.astype(np.float32)


def generate_description(feature: np.ndarray) -> str:
    """
    基于真实特征向量的确定性统计描述（占位实现）。
    说明: 真正医学级的面部特征描述（上唇上翘程度等）需要专门的
    面部特征回归/分类模型，这里是基于向量统计的确定性映射，
    保证 DeepSeek prompt 有真实数据可依。
    """
    if DESCRIPTION_MODE == "none":
        return ""

    # 取向量中固定区间的统计值映射到 6 个医学特征（0~1）
    segments = np.array_split(feature, 6)
    values = []
    for seg in segments:
        v = float(np.clip(np.mean(np.abs(seg)) * 3.0, 0.0, 1.0))
        values.append(round(v, 2))

    desc = {
        "上唇上翘程度": values[0],
        "鼻唇沟变浅程度": values[1],
        "下颌后缩程度": values[2],
        "牙列不齐程度": values[3],
        "张口呼吸特征": values[4],
        "面部比例失调程度": values[5],
        "特征向量能量": round(float(np.sum(feature ** 2)), 4),
    }
    return json.dumps(desc, ensure_ascii=False)

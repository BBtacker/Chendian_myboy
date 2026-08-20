"""
导出 EfficientNet-B3 特征提取 ONNX 模型
======================================
用法: python export_model.py
输出: models/efficientnet_b3.onnx （去掉分类头，输出 1536 维特征向量）

依赖: pip install torch torchvision onnx onnxscript onnxruntime
（onnx/onnxscript 供导出；onnxruntime 供末尾维度校验与运行时推理）
"""

import os

import torch
import torchvision.models as models

OUTPUT_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "models")
OUTPUT_PATH = os.path.join(OUTPUT_DIR, "efficientnet_b3.onnx")
INPUT_SIZE = 300  # EfficientNet-B3 标准输入


def main():
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    print("[1/3] 加载 EfficientNet-B3 预训练权重 (IMAGENET1K_V1)...")
    model = models.efficientnet_b3(weights=models.EfficientNet_B3_Weights.IMAGENET1K_V1)

    # 去掉分类头，只保留 backbone（features 输出 1536 维特征）
    model = model.features

    # 追加 Global Average Pooling，把 (C, H, W) 展平为 1536 维向量
    model = torch.nn.Sequential(
        model,
        torch.nn.AdaptiveAvgPool2d(1),
        torch.nn.Flatten(),
    )

    model.eval()

    print("[2/3] 导出 ONNX...")
    dummy = torch.randn(1, 3, INPUT_SIZE, INPUT_SIZE)
    torch.onnx.export(
        model,
        dummy,
        OUTPUT_PATH,
        input_names=["input"],
        output_names=["feature"],
        dynamic_axes={"input": {0: "batch"}, "feature": {0: "batch"}},
        opset_version=17,
    )
    print(f"[3/3] 导出成功: {OUTPUT_PATH}")

    # 验证输出维度
    import onnxruntime as ort
    sess = ort.InferenceSession(OUTPUT_PATH, providers=["CPUExecutionProvider"])
    shape = sess.get_outputs()[0].shape
    print(f"输出维度: {shape} (应为 [batch, 1536])")


if __name__ == "__main__":
    main()

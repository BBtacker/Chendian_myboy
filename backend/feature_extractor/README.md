# 特征提取服务（EfficientNet-B3 向量化）

## 作用
把诊断流水线中的「模拟特征提取」替换为真实 EfficientNet-B3 推理：
- 输入：面部照片
- 输出：1536 维特征向量（L2 归一化）+ 特征描述
- 特征向量存入 Milvus 用于 RAG 检索，特征描述拼入 DeepSeek prompt

## 你需要准备的数据（去网上找）

### 1. EfficientNet-B3 ONNX 模型（必需）
文件放入 `backend/feature_extractor/models/efficientnet_b3.onnx`

**获取方式（推荐，需要一台装了 Python 的机器）**：
```bash
pip install torch torchvision onnx
python export_model.py
```
`export_model.py` 会把 torchvision 的 EfficientNet-B3 预训练权重导出为 ONNX（去掉分类头，只保留 backbone，输出 1536 维特征）。

**或者**：搜索下载现成的 `efficientnet-b3.onnx` / `efficientnet_b3.onnx`（务必确认是**去掉分类头的特征模型**，输出维度 1536，而非 1000 类分类模型）。

> 判断方法：下载后运行 `python -c "import onnxruntime as ort; s=ort.InferenceSession('models/efficientnet_b3.onnx'); print(s.get_outputs()[0].shape)"`，输出应为 `[None, 1536]`（或 `[1, 1536]`）。

### 2.（可选）历史病例特征数据（RAG 冷启动）
系统运行后会自动积累（每次诊断把特征向量写入 Milvus）。如果想开箱就有检索效果，可准备一批标注了诊断结论的面部照片，调用 `/extract` 批量入库。**不准备也能跑**，只是前期 RAG 检索结果为空，DeepSeek 会退化为仅靠特征描述推理。

## 启动（本地开发）
```bash
cd backend/feature_extractor
pip install -r requirements.txt
bash run.sh
# 或: uvicorn main:app --host 0.0.0.0 --port 8085
```

## 测试
```bash
curl -X POST http://localhost:8085/extract -F "image=@某张照片.jpg"
# 返回 {"code":1, "feature_vector":[...1536个float...], "dimension":1536, "description":"{...}"}
```

## 与诊断服务对接
诊断服务的 `diagnosis.yaml` 增加：
```yaml
FeatureExtractor:
  URL: "http://127.0.0.1:8085"
  Timeout: 10   # 秒
```
诊断流水线会改为：读取图片 → 调用特征服务 → 真实向量入 Milvus + 检索 → DeepSeek 推理。

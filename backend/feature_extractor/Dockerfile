FROM python:3.11-slim

WORKDIR /app

# 安装依赖
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 复制代码
COPY app/ ./app/
COPY config.env ./

# 模型文件（构建时注入，或通过卷挂载）
# COPY models/efficientnet_b3.onnx ./models/efficientnet_b3.onnx

EXPOSE 8085

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8085"]

# 启动特征提取服务（本地开发）
uvicorn app.main:app --host 0.0.0.0 --port 8085 --reload

# 说明：
# 1. 先把 efficientnet_b3.onnx 放入 models/ 目录（见 README.md 获取方式）
# 2. 启动前确认依赖: pip install -r requirements.txt
# 3. 测试: curl -X POST http://localhost:8085/extract -F "image=@test.jpg"

// 调用后端API上传图片到阿里云OSS
export const uploadImageToOSS = async (file) => {
  try {
    const formData = new FormData();
    formData.append('file', file);

    console.log('上传图片到OSS:', file.name);
    const response = await fetch('/api/oss/upload', {
      method: 'POST',
      body: formData
    });

    console.log('OSS上传响应状态:', response.status, response.statusText);
    const result = await response.json();
    console.log('OSS上传响应内容:', result);

    if (result.success) {
      console.log('图片上传成功:', result.data);
      return result.data;
    } else {
      throw new Error(result.message || '图片上传失败');
    }
  } catch (error) {
    console.error('上传图片到OSS失败:', error);
    throw error;
  }
};

// 调用后端API从阿里云OSS删除图片
export const deleteImageFromOSS = async (imageUrl) => {
  try {
    console.log('从OSS删除图片:', imageUrl);
    const response = await fetch(`/api/oss/delete?imageUrl=${encodeURIComponent(imageUrl)}`, {
      method: 'DELETE'
    });

    console.log('OSS删除响应状态:', response.status, response.statusText);
    const result = await response.json();
    console.log('OSS删除响应内容:', result);

    if (result.success) {
      console.log('图片删除成功');
      return true;
    } else {
      throw new Error(result.message || '图片删除失败');
    }
  } catch (error) {
    console.error('从OSS删除图片失败:', error);
    throw error;
  }
};

// 将base64图片转换为文件对象
export const base64ToFile = (base64Data, filename) => {
  const arr = base64Data.split(',');
  const mime = arr[0].match(/:(.*?);/)[1];
  const bstr = atob(arr[1]);
  let n = bstr.length;
  const u8arr = new Uint8Array(n);
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }
  return new File([u8arr], filename, { type: mime });
};

// 调用后端API获取豆包回答
export const getDoubaoAnswer = async (messages) => {
  try {
    // 发起请求到后端API
    console.log('发起后端API请求:', messages);
    const response = await fetch('/api/doubao/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(messages)
    });

    console.log('后端API响应状态:', response.status, response.statusText);
    const result = await response.json();
    console.log('后端API响应内容:', result);

    // 检查响应是否成功
    if (result.success) {
      return result.answer;
    } else {
      console.error('后端API返回错误:', result.error);
      return `抱歉，我暂时无法回答您的问题。错误信息：${result.error}`;
    }
  } catch (error) {
    console.error('调用后端API失败:', error);
    return '抱歉，我暂时无法回答您的问题，请稍后再试。';
  }
};

// 调用后端API生成图片
export const generateImage = async (imageData, mood) => {
  try {
    // 发起请求到后端API
    console.log('发起图片生成API请求:', { mood, imageDataLength: imageData ? imageData.length : 0 });
    const response = await fetch('/api/doubao/generate-image', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ imageData, mood })
    });

    console.log('后端API响应状态:', response.status, response.statusText);
    
    // 增加对非 200 状态码的处理
    if (!response.ok) {
        throw new Error(`网络请求失败: ${response.status}`);
    }

    const result = await response.json();
    console.log('后端API响应内容:', result);
    
    // 检查响应是否成功
    if (result.success) {
      console.log('图片生成成功');
      // 检查是否是URL
      if (result.image && result.image.startsWith('http')) {
        console.log('返回的是URL:', result.image);
        return result.image;
      } else {
        console.log('返回的是base64数据');
        return result.image;
      }
    } else {
      // 直接抛出后端返回的具体错误信息
      console.error('后端API返回业务错误:', result.error);
      throw new Error(result.error || '图片生成失败，原因未知');
    }
  } catch (error) {
    console.error('调用后端API发生异常:', error);
    // 将错误抛给调用者（DreamPage.jsx），让它去弹窗提示
    throw error;
  }
};

// 聊天历史记录管理
export const chatHistoryManager = {
  // 获取当前用户的聊天历史记录
  getChatHistory: () => {
    const currentUser = localStorage.getItem('currentUser');
    if (!currentUser) return [];
    
    const user = JSON.parse(currentUser);
    const chatHistoryKey = `chatHistory_${user.id}`;
    const chatHistory = localStorage.getItem(chatHistoryKey);
    
    return chatHistory ? JSON.parse(chatHistory) : [];
  },
  
  // 保存聊天历史记录
  saveChatHistory: (messages) => {
    const currentUser = localStorage.getItem('currentUser');
    if (!currentUser) return false;
    
    const user = JSON.parse(currentUser);
    const chatHistoryKey = `chatHistory_${user.id}`;
    localStorage.setItem(chatHistoryKey, JSON.stringify(messages));
    return true;
  },
  
  // 清除聊天历史记录
  clearChatHistory: () => {
    const currentUser = localStorage.getItem('currentUser');
    if (!currentUser) return false;
    
    const user = JSON.parse(currentUser);
    const chatHistoryKey = `chatHistory_${user.id}`;
    localStorage.removeItem(chatHistoryKey);
    return true;
  }
};
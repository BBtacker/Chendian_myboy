// API请求工具
//const API_BASE_URL = 'http://localhost:8080/api';
const API_BASE_URL = '/api';

// 通用请求方法
async function request(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const defaultOptions = {
    credentials: 'include',
  };
  
  // 如果请求体是FormData类型，不设置默认的Content-Type
  if (!(options.body instanceof FormData)) {
    defaultOptions.headers = {
      'Content-Type': 'application/json',
    };
  }
  
  const mergedOptions = {
    ...defaultOptions,
    ...options,
    headers: options.headers === undefined 
      ? defaultOptions.headers 
      : {
          ...(defaultOptions.headers || {}),
          ...options.headers
        },
  };
  
  try {
    const response = await fetch(url, mergedOptions);
    
    // 检查响应是否为空
    const contentType = response.headers.get('content-type');
    let responseData;
    
    if (contentType && contentType.includes('application/json')) {
      responseData = await response.json();
    } else {
      responseData = await response.text();
    }
    
    if (!response.ok) {
      // 创建错误对象，包含响应状态和数据
      const error = new Error(`HTTP error! status: ${response.status}`);
      error.status = response.status;
      error.responseData = responseData;
      throw error;
    }
    
    return responseData;
  } catch (error) {
    console.error('API请求错误:', error);
    throw error;
  }
}

// 人脸分析相关API
export const analysisApi = {
  // 获取所有分析记录
  getAllRecords: () => request('/analysis/records'),
  
  // 获取最新的分析记录
  getLatestRecords: (limit = 10) => request(`/analysis/records/latest?limit=${limit}`),
  
  // 根据ID获取分析记录
  getRecordById: (id) => request(`/analysis/records/${id}`),
  
  // 根据日期范围获取分析记录
  getRecordsByDateRange: (startDate, endDate) => {
    // 格式化日期为'yyyy-MM-dd HH:mm:ss'
    const formatDate = (date) => {
      return date.toISOString().slice(0, 19).replace('T', ' ');
    };
    
    const formattedStart = formatDate(startDate);
    const formattedEnd = formatDate(endDate);
    
    return request(`/analysis/records/date-range?startDate=${formattedStart}&endDate=${formattedEnd}`);
  },
  
  // 根据主要表情获取分析记录
  getRecordsByDominantExpression: (expression) => {
    return request(`/analysis/records/dominant-expression/${expression}`);
  },
  
  // 上传图片进行分析
  analyzeImage: (imageFile) => {
    const formData = new FormData();
    formData.append('file', imageFile);
    
    return request('/analysis/analyze', {
      method: 'POST',
      body: formData,
      headers: undefined, // 不设置headers，完全使用默认值
    });
  },
  
  // 删除分析记录
  deleteRecord: (id) => request(`/analysis/records/${id}`, {
    method: 'DELETE',
  }),
  
  // 更新分析记录的表情数据（手动录入）
  updateEmotions: (id, emotionData) => request(`/analysis/records/${id}/emotions`, {
    method: 'PUT',
    body: JSON.stringify(emotionData),
  }),
};

// 豆包API请求
export const doubaoApi = {
  // 发送消息到豆包API
  sendMessage: (message) => {
    // 这里是模拟豆包API的实现
    // 在实际项目中，你需要替换为真实的豆包API调用
    return new Promise((resolve) => {
      setTimeout(() => {
        // 模拟API响应
        resolve({
          response: getDoubaoResponse(message),
        });
      }, 1000);
    });
  },
};

// 用户认证相关API
export const authApi = {
  // 用户登录
  login: (username, password) => request('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  }),
  
  // 用户注册
  register: (username, password) => request('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  }),
  
  // 获取用户信息
  getUserById: (id) => request(`/auth/user/${id}`),
};

// 日记相关API
export const diaryApi = {
  // 创建日记
  createDiary: (diary) => request('/diary', {
    method: 'POST',
    body: JSON.stringify(diary),
  }),
  
  // 根据用户ID获取所有日记
  getDiariesByUserId: (userId) => {
    if (!userId) {
      throw new Error('用户ID不能为空');
    }
    return request(`/diary/user/${userId}`);
  },
  
  // 根据ID获取日记
  getDiaryById: (id, userId) => {
    if (!userId) {
      throw new Error('用户ID不能为空');
    }
    return request(`/diary/${id}/${userId}`);
  },
  
  // 根据用户ID和日期范围获取日记
  getDiariesByUserIdAndDateRange: (userId, startDate, endDate) => {
    if (!userId) {
      throw new Error('用户ID不能为空');
    }
    return request(`/diary/user/${userId}/date-range?startDate=${startDate}&endDate=${endDate}`);
  },
  
  // 更新日记
  updateDiary: (diary) => request(`/diary/${diary.id}`, {
    method: 'PUT',
    body: JSON.stringify(diary),
  }),
  
  // 删除日记
  deleteDiary: (id, userId) => {
    if (!userId) {
      throw new Error('用户ID不能为空');
    }
    return request(`/diary/${id}/${userId}`, {
      method: 'DELETE',
    });
  },
  
  // 为日记添加照片
  addDiaryPhoto: (photo) => request(`/diary/${photo.diaryId}/photos`, {
    method: 'POST',
    body: JSON.stringify(photo),
  }),
  
  // 获取日记的照片列表
  getPhotosByDiaryId: (diaryId) => request(`/diary/${diaryId}/photos`),
};

// 模拟豆包API响应
function getDoubaoResponse(message) {
  const lowerMessage = message.toLowerCase();
  
  if (lowerMessage.includes('表情识别')) {
    return '表情识别是一种通过分析人脸图像来识别人类表情的技术。它可以识别多种基本表情，如开心、悲伤、愤怒、惊讶、恐惧、厌恶和中性等。';
  } else if (lowerMessage.includes('如何使用')) {
    return '您可以点击"选择图片"按钮上传一张包含人脸的图片，然后点击"开始检测"按钮，系统会自动分析图片中的人脸表情。';
  } else if (lowerMessage.includes('支持哪些表情')) {
    return '我们的系统支持识别7种基本表情：开心、悲伤、愤怒、惊讶、恐惧、厌恶和中性。';
  } else if (lowerMessage.includes('准确率')) {
    return '我们的表情识别系统在标准数据集上的准确率约为92%，实际应用中可能会因图像质量、光照条件等因素有所变化。';
  } else if (lowerMessage.includes('历史记录')) {
    return '您可以在"历史人脸检测查询"页面查看所有的历史检测记录。每条记录包含检测日期、时间、主要表情以及各表情的分布情况。';
  } else {
    return '感谢您的提问！我是专门为表情识别系统设计的智能助手，您可以问我关于表情识别的相关问题。';
  }
}

// ===== 用户个人信息相关API =====
export const userApi = {
  updateUser: (user) => request('/auth/user/update', {
    method: 'PUT',
    body: JSON.stringify(user),
  }),
  getUserById: (id) => request(`/auth/user/${id}`),
};

// ===== 游戏得分相关API =====
export const gameApi = {
  saveScore: (scoreData) => request('/game/score', {
    method: 'POST',
    body: JSON.stringify(scoreData),
  }),
  getScores: (userId) => request(`/game/scores/${userId}`),
  getBestScore: (userId) => request(`/game/best/${userId}`),
  getHighScore: (userId) => request(`/game/high-score/${userId}`),
};

// ===== 绘梦记录相关API =====
export const dreamApi = {
  saveDreamImage: (data) => request('/dream-image/save', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  getDreamImages: (userId) => request(`/dream-image/list/${userId}`),
  deleteDreamImage: (id, userId) => request(`/dream-image/${id}/user/${userId}`, {
    method: 'DELETE',
  }),
};

// ===== AI聊天相关API =====
export const chatApi = {
  createSession: (userId, sessionName) => request('/chat/session/create', {
    method: 'POST',
    body: JSON.stringify({ userId, sessionName }),
  }),
  getSessions: (userId) => request(`/chat/sessions/${userId}`),
  getMessages: (sessionId) => request(`/chat/messages/${sessionId}`),
  saveMessage: (message) => request('/chat/message/save', {
    method: 'POST',
    body: JSON.stringify(message),
  }),
  deleteSession: (sessionId) => request(`/chat/session/${sessionId}`, {
    method: 'DELETE',
  }),
};

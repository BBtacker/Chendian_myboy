// 日记数据管理工具

// 从localStorage获取日记数据
export const getDiariesFromStorage = () => {
  const diaries = localStorage.getItem('diaries');
  return diaries ? JSON.parse(diaries) : [];
};

// 检查localStorage剩余空间
const checkLocalStorageQuota = (data) => {
  try {
    // 创建一个临时条目来测试存储配额
    const testKey = '__test_quota__';
    const testValue = JSON.stringify(data);
    localStorage.setItem(testKey, testValue);
    localStorage.removeItem(testKey);
    return true;
  } catch (e) {
    return false;
  }
};

// 保存日记数据到localStorage
export const saveDiariesToStorage = (diaries) => {
  const data = JSON.stringify(diaries);
  
  // 检查localStorage是否有足够空间
  if (checkLocalStorageQuota(diaries)) {
    localStorage.setItem('diaries', data);
    return true;
  } else {
    console.error('localStorage配额不足，无法保存日记数据');
    // 可以考虑添加清理旧数据的逻辑
    return false;
  }
};

// 添加新日记
export const addDiary = (diary) => {
  const diaries = getDiariesFromStorage();
  const newDiary = {
    id: Date.now(),
    ...diary,
    createdAt: new Date().toISOString()
  };
  const updatedDiaries = [...diaries, newDiary];
  
  const success = saveDiariesToStorage(updatedDiaries);
  if (success) {
    return newDiary;
  } else {
    throw new Error('保存日记失败：localStorage配额不足');
  }
};

// 删除日记
export const deleteDiary = (diaryId) => {
  const diaries = getDiariesFromStorage();
  const updatedDiaries = diaries.filter(diary => diary.id !== diaryId);
  
  const success = saveDiariesToStorage(updatedDiaries);
  if (success) {
    return updatedDiaries;
  } else {
    throw new Error('删除日记失败：localStorage配额不足');
  }
};

// 筛选日记
export const filterDiaries = (diaries, filters) => {
  let filteredDiaries = [...diaries];
  
  // 按日期范围筛选
  if (filters.startDate && filters.endDate) {
    filteredDiaries = filteredDiaries.filter(diary => {
      // 使用selectedDate字段，因为保存时使用的是这个字段名
      const diaryDate = new Date(diary.selectedDate);
      const startDate = new Date(filters.startDate);
      const endDate = new Date(filters.endDate);
      return diaryDate >= startDate && diaryDate <= endDate;
    });
  }
  
  // 按心情筛选
  if (filters.mood) {
    // 使用selectedMood字段，因为保存时使用的是这个字段名
    filteredDiaries = filteredDiaries.filter(diary => diary.selectedMood === filters.mood);
  }
  
  return filteredDiaries;
};

// 从localStorage获取用户数据
export const getUserFromStorage = () => {
  const user = localStorage.getItem('currentUser');
  return user ? JSON.parse(user) : {
    id: 1,
    username: '用户123',
    avatar: 'https://picsum.photos/seed/user1/200/200',
    coverPhoto: '',
    bio: '天天开心~',
    followers: 120,
    following: 80,
    likes: 450
  };
};

// 保存用户数据到localStorage
export const saveUserToStorage = (user) => {
  try {
    localStorage.setItem('currentUser', JSON.stringify(user));
    return true;
  } catch (e) {
    console.error('localStorage配额不足，无法保存用户数据');
    return false;
  }
};

// 评论数据管理

// 从localStorage获取评论数据
export const getCommentsFromStorage = () => {
  const comments = localStorage.getItem('comments');
  return comments ? JSON.parse(comments) : [];
};

// 保存评论数据到localStorage
export const saveCommentsToStorage = (comments) => {
  try {
    localStorage.setItem('comments', JSON.stringify(comments));
    return true;
  } catch (e) {
    console.error('localStorage配额不足，无法保存评论数据');
    return false;
  }
};

// 添加新评论
export const addComment = (comment) => {
  const comments = getCommentsFromStorage();
  const newComment = {
    id: Date.now(),
    ...comment,
    createdAt: new Date().toISOString()
  };
  const updatedComments = [...comments, newComment];
  
  const success = saveCommentsToStorage(updatedComments);
  if (success) {
    return newComment;
  } else {
    throw new Error('保存评论失败：localStorage配额不足');
  }
};

// 删除评论
export const deleteComment = (commentId) => {
  const comments = getCommentsFromStorage();
  const updatedComments = comments.filter(comment => comment.id !== commentId);
  
  const success = saveCommentsToStorage(updatedComments);
  if (success) {
    return updatedComments;
  } else {
    throw new Error('删除评论失败：localStorage配额不足');
  }
};

// 获取特定日记的评论
export const getCommentsByDiaryId = (diaryId) => {
  const comments = getCommentsFromStorage();
  return comments.filter(comment => comment.diaryId === diaryId);
};
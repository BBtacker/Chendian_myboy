import axios from 'axios'

// 创建axios实例
const request = axios.create({
  baseURL: '/api', // 使用/api作为基础路径，通过vite代理转发到后端
  timeout: 30000, // 增加超时时间到30秒，因为图片分析可能需要较长时间
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 从localStorage中获取token
    const token = localStorage.getItem('token')
    if (token) {
      // 在请求头中添加Authorization字段
      config.headers.Authorization = token
    }
    
    // 对于文件上传请求，删除Content-Type，让浏览器自动设置
    if (config.data instanceof FormData) {
      delete config.headers['Content-Type'];
    }
    
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    // 直接返回响应数据
    return response.data
  },
  error => {
    return Promise.reject(error)
  }
)

// 添加一个专门用于面部检测的实例，设置更长的超时时间
const faceAnalysisRequest = axios.create({
  baseURL: '/api',
  timeout: 120000, // 设置超时时间为2分钟，因为面部检测可能需要较长时间
  headers: {
    'Content-Type': 'multipart/form-data'
  }
})

// 面部检测请求拦截器
faceAnalysisRequest.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = token
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 面部检测响应拦截器
faceAnalysisRequest.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    return Promise.reject(error)
  }
)

// 添加一个专门用于下载的实例
const downloadRequest = axios.create({
  baseURL: '/api',
  timeout: 30000,
  responseType: 'blob' // 设置响应类型为blob，用于文件下载
})

// 下载请求拦截器
downloadRequest.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = token
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 下载响应拦截器
downloadRequest.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          localStorage.removeItem('token')
          window.location.href = '/'
          break
        default:
          console.error(`下载错误${error.response.status}`)
      }
    } else {
      console.error('网络错误')
    }
    return Promise.reject(error)
  }
)

export { downloadRequest, faceAnalysisRequest }
export default request
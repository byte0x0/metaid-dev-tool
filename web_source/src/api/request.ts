import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'
import { ElMessage } from 'element-plus'

// 扩展 Error 类型
interface CustomError extends Error {
  response?: AxiosResponse
}

// 定义响应数据类型
interface ApiResponse {
  code: number
  msg?: string
  message?: string
  data?: any
}

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 在这里可以添加认证信息等
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { data } = response
    // 如果是代理接口，直接返回原始数据
    if (response.config.url?.includes('/proxy/utxo')) {
      return data
    }
    // 处理成功状态码
    if (data.code === 0 || data.code === 2000) {
      return data.data || data.msg || data
    }
    // 保留原始错误信息
    const error = new Error(data.msg || data.message || '请求失败') as CustomError
    error.response = response
    return Promise.reject(error)
  },
  (error: AxiosError<ApiResponse>) => {
    // 保留原始错误信息
    if (error.response?.data?.msg) {
      error.message = error.response.data.msg
    } else if (error.response?.data?.message) {
      error.message = error.response.data.message
    }
    return Promise.reject(error)
  }
)

// 封装 GET 请求
export const get = <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
  return request.get(url, config)
}

// 封装 POST 请求
export const post = <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
  return request.post(url, data, config)
}

// 封装 PUT 请求
export const put = <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
  return request.put(url, data, config)
}

// 封装 DELETE 请求
export const del = <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
  return request.delete(url, config)
}

export default request 
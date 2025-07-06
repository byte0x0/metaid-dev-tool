// API响应基础接口
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页请求参数接口
export interface PaginationParams {
  page: number
  pageSize: number
}

// 分页响应数据接口
export interface PaginationData<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 通用响应类型
export type ApiResult<T> = Promise<ApiResponse<T>>

// 分页响应类型
export type PaginationResult<T> = Promise<ApiResponse<PaginationData<T>>> 
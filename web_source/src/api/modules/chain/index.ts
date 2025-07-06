import { get, post, put, del } from '../../request'
import type {
  ChainInfo,
  CreateChainParams,
  UpdateChainParams,
  ChainListParams,
  ChainListResponse,
  ChainDetailResponse,
} from './types'

// 获取Chain列表
export function getChainList(params: ChainListParams): Promise<ChainListResponse> {
  return get<ChainListResponse>('/chains', { params })
}

// 获取Chain详情
export function getChainDetail(id: string): Promise<ChainInfo> {
  return get<ChainDetailResponse>(`/chains/${id}`)
}

// 创建Chain
export function createChain(data: CreateChainParams): Promise<ChainInfo> {
  return post<ChainDetailResponse>('/chains', data)
}

// 更新Chain
export function updateChain(id: string, data: UpdateChainParams): Promise<ChainInfo> {
  return put<ChainDetailResponse>(`/chains/${id}`, data)
}

// 删除Chain
export function deleteChain(id: string): Promise<void> {
  return del(`/chains/${id}`)
}


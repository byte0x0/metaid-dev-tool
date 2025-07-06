import type { PaginationParams, PaginationData } from '../../types'

// Chain基本信息接口
export interface ChainInfo {
  ID: number
  name: string
  broadcast_url: string
  utxo_url: string
  chain_type: string
  CreatedAt: string
  UpdatedAt: string
}

// Chain创建参数接口
export interface CreateChainParams {
  name: string
  broadcastUrl: string
  utxoUrl: string
  chainType: string
}

// Chain更新参数接口
export interface UpdateChainParams {
  name?: string
  broadcastUrl?: string
  utxoUrl?: string
  chainType?: string
}

// Chain列表查询参数
export interface ChainListParams {
  chainType?: string
}

// Chain列表响应类型
export type ChainListResponse = ChainInfo[]

// Chain详情响应类型
export type ChainDetailResponse = ChainInfo 
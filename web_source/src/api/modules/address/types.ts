import type { PaginationParams } from '../../types'

// 地址类型枚举
export enum AddressType {
  Taproot = 'taproot',
  Segwit = 'segwit'
}

// 地址信息接口
export interface AddressInfo {
  ID: number
  address: string
  private_key: string
  type: AddressType
  chain_id: number
  CreatedAt: string
  UpdatedAt: string
  Chain?: {
    ID: number
    name: string
    chain_type: string
  }
}

// 创建地址参数接口
export interface CreateAddressParams {
  chain_id: number
  type: AddressType
}

// 地址列表查询参数
export interface AddressListParams extends PaginationParams {
  type?: AddressType
  chain_id?: number
}

// 地址列表响应类型
export type AddressListResponse = AddressInfo[]

// 地址详情响应类型
export type AddressDetailResponse = AddressInfo 
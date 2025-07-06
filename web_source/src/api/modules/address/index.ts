import { get, post, del } from '../../request'
import type {
  AddressInfo,
  AddressListParams,
  AddressListResponse,
  AddressDetailResponse,
  CreateAddressParams
} from './types'

// 获取地址列表
export function getAddressList(params: AddressListParams): Promise<AddressListResponse> {
  return get<AddressListResponse>('/addresses', { params })
}

// 获取地址详情
export function getAddressDetail(id: number): Promise<AddressInfo> {
  return get<AddressDetailResponse>(`/addresses/${id}`)
}

// 创建地址
export function createAddress(data: CreateAddressParams): Promise<AddressInfo> {
  return post<AddressInfo>('/addresses', data)
}

// 删除地址
export function deleteAddress(id: number): Promise<void> {
  return del(`/addresses/${id}`)
} 
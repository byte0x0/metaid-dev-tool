import { get, post } from '../../request'
import type {
  CreateInscribeRequest,
  InscribeResponse,
  InscribeInfo,
  InscribeListParams,
  InscribeListResponse,
  Utxo
} from './types'

// 创建铭文
export const createInscribe = (data: CreateInscribeRequest) => {
  return post<InscribeResponse>('/inscribes', data)
}

// 获取铭文详情
export const getInscribe = (id: number) => {
  return get<InscribeInfo>(`/inscribe/${id}`)
}

// 获取铭文列表
export const listInscribes = (params: InscribeListParams) => {
  return get<InscribeListResponse>('/inscribe', { params })
}

// 获取 UTXO 列表
export const getUtxoList = (address: string) => {
  return get<Utxo[]>('/utxo', { params: { address } })
}

// 广播交易
export const broadcastTx = (txRaw: string) => {
  return post<{ txId: string }>('/broadcast', { txRaw })
} 
export interface CreateInscribeRequest {
  chainId: number
  metaIdFlag: string
  path: string
  operation: string
  payload: any
  pinOutValue: number
  address: string
  feeRate: number
  utxoList: Utxo[]
  commitTx?: string
  revealTx?: string
}

export interface Utxo {
  outTxId: string
  outIndex: number
  outValue: number
  outRaw?: string
}

export interface InscribeResponse {
  commit_tx_raw: string
  reveal_tx_raw: string
  commit_tx_id: string
  miner_fee: number
}

export interface InscribeInfo {
  id: number
  chainId: number
  status: string
  commitTxId: string
  revealTxId: string
  createdAt: string
  updatedAt: string
}

export interface InscribeListParams {
  status?: string
  chainId?: number
}

export interface InscribeListResponse {
  list: InscribeInfo[]
  total: number
} 
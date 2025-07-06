<template>
  <div class="pin-container">
    <el-form :model="form" label-width="120px" :rules="rules" ref="formRef">
      <el-form-item label="Chain" prop="chainId">
        <el-select 
          v-model="form.chainId" 
          placeholder="Select chain"
          @change="handleChainChange"
          :clearable="true"
        >
          <el-option
            v-for="chain in chainList"
            :key="chain.ID"
            :label="chain.name"
            :value="chain.ID"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="MetaID Flag" prop="metaIdFlag">
        <el-input v-model="form.metaIdFlag" placeholder="Enter MetaID flag" />
      </el-form-item>

      <el-form-item label="Path" prop="path">
        <el-input v-model="form.path" placeholder="Enter path" />
      </el-form-item>

      <el-form-item label="Operation" prop="operation">
        <el-select v-model="form.operation" placeholder="Select operation">
          <el-option label="Create" value="create" />
          <el-option label="Modify" value="modify" />
          <el-option label="Revoke" value="revoke" />
          <el-option label="Hide" value="hide" />
        </el-select>
      </el-form-item>

      <el-form-item label="Payload" prop="payload">
        <el-input
          v-model="form.payload"
          type="textarea"
          :rows="4"
          placeholder="Enter payload"
        />
      </el-form-item>

      <el-form-item label="PIN Out Value" prop="pinOutValue">
        <el-input-number v-model="form.pinOutValue" :min="0" />
      </el-form-item>

      <el-form-item label="Address" prop="address">
        <div class="address-input">
          <el-select 
            v-model="form.address" 
            placeholder="Select address"
            :disabled="!form.chainId"
            @change="(val: string) => console.log('Selected address:', val)"
          >
            <el-option
              v-for="address in addressList"
              :key="address.ID"
              :label="address.address"
              :value="address.address"
            />
          </el-select>
          <el-button 
            type="primary" 
            @click="() => {
              try {
                const chain = getCurrentChain()
                if (chain?.utxo_url) {
                  showUtxoDrawer(form.address, chain.utxo_url)
                } else {
                  ElMessage.warning('UTXO URL not configured for this chain')
                }
              } catch (error) {
                console.error('Error:', error)
                ElMessage.error('Failed to get UTXO URL')
              }
            }"
          >选择 UTXO</el-button>
        </div>
      </el-form-item>

      <el-form-item label="Fee Rate" prop="feeRate">
        <el-input-number v-model="form.feeRate" :min="1" />
      </el-form-item>

      <!-- 已选择的 UTXO 列表 -->
      <el-form-item v-if="form.utxoList.length > 0" label="Selected UTXOs">
        <el-table :data="form.utxoList" style="width: 100%">
          <el-table-column prop="outTxId" label="Transaction ID" />
          <el-table-column prop="outIndex" label="Output Index" />
          <el-table-column prop="outValue" label="Value">
            <template #default="{ row }">
              {{ row.outValue / 100000000 }} BTC
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="Operations" width="120">
            <template #default="{ $index }">
              <el-button type="danger" link @click="removeUtxo($index)">Remove</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleSubmit">Create PIN</el-button>
        <el-button 
          type="success" 
          @click="handleBroadcast('commit')"
          :disabled="!form.commitTx"
        >广播 Commit 交易</el-button>
        <el-button 
          type="warning" 
          @click="handleBroadcast('reveal')"
          :disabled="!form.revealTx"
        >广播 Reveal 交易</el-button>
      </el-form-item>

      <!-- 广播结果显示 -->
      <el-form-item v-if="broadcastResult.commit || broadcastResult.reveal">
        <div class="broadcast-result">
          <h4>广播结果</h4>
          <el-descriptions :column="1" border>
            <el-descriptions-item v-if="broadcastResult.commit">
              <template #label>
                <div class="result-label">Commit Transaction ID</div>
              </template>
              <div class="tx-content">
                <el-input v-model="broadcastResult.commit" type="text" readonly />
                <el-button type="primary" link @click="copyContent(broadcastResult.commit)">
                  <el-icon><Document /></el-icon>
                  复制
                </el-button>
              </div>
            </el-descriptions-item>
            <el-descriptions-item v-if="broadcastResult.reveal">
              <template #label>
                <div class="result-label">Reveal Transaction ID</div>
              </template>
              <div class="tx-content">
                <el-input v-model="broadcastResult.reveal" type="text" readonly />
                <el-button type="primary" link @click="copyContent(broadcastResult.reveal)">
                  <el-icon><Document /></el-icon>
                  复制
                </el-button>
              </div>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, inject, onMounted, computed, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { createInscribe } from '../api/modules/inscribe'
import { getChainList } from '../api/modules/chain'
import { getAddressList } from '../api/modules/address'
import type { CreateInscribeRequest, Utxo } from '../api/modules/inscribe/types'
import type { ChainInfo, ChainListResponse } from '../api/modules/chain/types'
import type { AddressInfo } from '../api/modules/address/types'
import { post } from '../api/request'
import { Document } from '@element-plus/icons-vue'

const formRef = ref<FormInstance>()
const chainList = ref<ChainListResponse>([])
const addressList = ref<AddressInfo[]>([])
const currentChain = ref<ChainInfo | null>(null)

const form = reactive<CreateInscribeRequest & { commitTx?: string; revealTx?: string }>({
  chainId: null as unknown as number,
  metaIdFlag: 'metaid',
  path: '',
  operation: 'create',
  payload: '',
  pinOutValue: 546,
  address: '',
  feeRate: 1,
  utxoList: [],
  commitTx: '',
  revealTx: ''
})

const rules = reactive<FormRules>({
  chainId: [{ required: true, message: 'Please select chain', trigger: 'change' }],
  metaIdFlag: [{ required: true, message: 'Please enter MetaID flag', trigger: 'blur' }],
  path: [{ required: true, message: 'Please enter path', trigger: 'blur' }],
  operation: [{ required: true, message: 'Please select operation', trigger: 'change' }],
  payload: [{ required: true, message: 'Please enter payload', trigger: 'blur' }],
  pinOutValue: [{ required: true, message: 'Please enter PIN out value', trigger: 'blur' }],
  address: [{ required: true, message: 'Please select address', trigger: 'change' }],
  feeRate: [{ required: true, message: 'Please enter fee rate', trigger: 'blur' }]
})

// 注入父组件提供的方法
const showUtxoDrawer = inject('showUtxoDrawer') as (address: string, utxoUrl: string) => void
const setCreateResult = inject('setCreateResult') as (result: any) => void

// 获取当前选中的链
const selectedChain = computed(() => {
  return chainList.value?.find(c => c.ID === form.chainId)
})

// 监听 chainId 变化
watch(() => form.chainId, (newChainId) => {
  if (newChainId) {
    currentChain.value = chainList.value.find(c => c.ID === newChainId) || null
  } else {
    currentChain.value = null
  }
})

// 获取链列表
const fetchChainList = async () => {
  try {
    const response = await getChainList({})
    chainList.value = response
    console.log('Chain list:', response)
  } catch (error) {
    ElMessage.error('Failed to fetch chain list')
    console.error('Error:', error)
  }
}

// 获取地址列表
const fetchAddressList = async (chainId: number) => {
  try {
    const response = await getAddressList({ 
      chain_id: chainId,
      page: 1,
      pageSize: 100
    })
    console.log('Address list:', response)
    addressList.value = response
    console.log('Updated addressList:', addressList.value)
  } catch (error) {
    ElMessage.error('Failed to fetch address list')
    console.error('Error:', error)
  }
}

// 处理链选择变化
const handleChainChange = async (chainId: number) => {
  console.log('Chain changed:', chainId)
  form.chainId = chainId
  form.address = '' // 清空地址选择
  if (chainId) {
    await fetchAddressList(chainId)
  } else {
    addressList.value = []
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const response = await createInscribe(form)
        form.commitTx = response.commit_tx_raw
        form.revealTx = response.reveal_tx_raw
        ElMessage.success('PIN created successfully')
        setCreateResult(response)
      } catch (error) {
        ElMessage.error('Failed to create PIN')
        console.error('Error:', error)
      }
    }
  })
}

const addUtxo = (utxo: Utxo) => {
  form.utxoList.push(utxo)
}

const removeUtxo = (index: number) => {
  form.utxoList.splice(index, 1)
}

// 获取当前链
const getCurrentChain = () => {
  return chainList.value.find(c => c.ID === form.chainId)
}

// 广播结果
const broadcastResult = reactive({
  commit: '',
  reveal: ''
})

// 复制内容到剪贴板
const copyContent = async (content: string) => {
  try {
    await navigator.clipboard.writeText(content)
    ElMessage.success('复制成功')
  } catch (error) {
    ElMessage.error('复制失败')
    console.error('Error:', error)
  }
}

// 广播交易
const handleBroadcast = async (type: 'commit' | 'reveal') => {
  try {
    const chain = getCurrentChain()
    if (!chain?.broadcast_url) {
      ElMessage.warning('Broadcast URL not configured for this chain')
      return
    }

    const txContent = type === 'commit' ? form.commitTx : form.revealTx
    if (!txContent) {
      ElMessage.warning(`No ${type} transaction to broadcast`)
      return
    }

    const response = await post('/broadcast', {
      url: chain.broadcast_url,
      content: txContent
    })

    // 更新广播结果
    if (type === 'commit') {
      broadcastResult.commit = response
    } else {
      broadcastResult.reveal = response
    }

    ElMessage.success(`${type} transaction broadcasted successfully`)
    console.log('Broadcast response:', response)
  } catch (error: any) {
    // 从错误响应中提取具体错误信息
    let errorMsg = 'Failed to broadcast transaction'
    if (error.response?.data?.msg) {
      errorMsg = error.response.data.msg
    } else if (error.response?.data?.message) {
      errorMsg = error.response.data.message
    } else if (error.message) {
      errorMsg = error.message
    }
    
    // 只显示一个错误提示
    ElMessage({
      type: 'error',
      message: `${type} 交易广播失败: ${errorMsg}`,
      duration: 5000
    })
    console.error('Broadcast error:', error)
  }
}

// 初始化
onMounted(() => {
  fetchChainList()
})

// 暴露方法给父组件调用
defineExpose({
  addUtxo,
  setCurrentAddress: (address: string) => {
    form.address = address
  }
})
</script>

<style scoped>
.pin-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.address-input {
  display: flex;
  gap: 10px;
  width: 100%;
}

.address-input .el-select {
  flex: 1;
  min-width: 400px;
}

.address-input .el-button {
  flex-shrink: 0;
}

.broadcast-result {
  margin-top: 20px;
  padding: 20px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
}

.broadcast-result h4 {
  margin-bottom: 20px;
  color: var(--el-text-color-primary);
  font-size: 16px;
  font-weight: 600;
}

.result-label {
  font-weight: 500;
  color: var(--el-text-color-regular);
  margin-bottom: 8px;
}

.tx-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.tx-content .el-input {
  flex: 1;
}

:deep(.el-descriptions__cell) {
  padding: 16px;
}

:deep(.el-descriptions__label) {
  width: 100%;
  padding-bottom: 8px;
  background-color: var(--el-bg-color);
}

:deep(.el-descriptions__content) {
  padding: 0;
}

:deep(.el-descriptions__body) {
  background-color: var(--el-bg-color);
}

:deep(.el-input__wrapper) {
  box-shadow: none;
  background-color: var(--el-fill-color-light);
}

:deep(.el-textarea__inner) {
  font-family: monospace;
  font-size: 14px;
  line-height: 1.5;
}
</style>

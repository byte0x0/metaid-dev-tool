<template>
  <div class="utxo-container">
    <el-table 
      v-if="utxoList.length > 0" 
      :data="utxoList" 
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="outTxId" label="Transaction ID" />
      <el-table-column prop="outIndex" label="Output Index" />
      <el-table-column prop="outValue" label="Value">
        <template #default="{ row }">
          {{ row.outValue / 100000000 }} BTC
        </template>
      </el-table-column>
    </el-table>

    <div v-else class="empty-state">
      <el-empty description="No UTXOs found" />
    </div>

    <div class="action-buttons">
      <el-button type="primary" @click="handleConfirm">确认选择</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { Utxo } from '../api/modules/inscribe/types'
import { post } from '../api/request'

const utxoList = ref<Utxo[]>([])
const selectedUtxos = ref<Utxo[]>([])

const handleGetUtxo = async (address: string) => {
  try {
    if (!props.utxoUrl) {
      ElMessage.error('UTXO URL not configured')
      return
    }
    const response = await post('/proxy/utxo', {
      url: props.utxoUrl.replace('{address}', address)
    })
    // 转换数据结构
    utxoList.value = response.result.map((item: any) => ({
      outTxId: item.txId,
      outIndex: item.vout,
      outValue: item.satoshi
    }))
    if (utxoList.value.length === 0) {
      ElMessage.warning('No UTXOs found for this address')
    }
  } catch (error) {
    ElMessage.error('Failed to get UTXO list')
    console.error('Error:', error)
  }
}

const handleSelectionChange = (selection: Utxo[]) => {
  selectedUtxos.value = selection
}

const handleConfirm = () => {
  if (selectedUtxos.value.length === 0) {
    ElMessage.warning('请选择至少一个 UTXO')
    return
  }
  // 触发选择事件，将选中的 UTXOs 传递给父组件
  emit('select', selectedUtxos.value)
}

const emit = defineEmits<{
  (e: 'select', utxos: Utxo[]): void
}>()

// 接收父组件传入的地址和 UTXO URL
const props = defineProps<{
  address?: string
  utxoUrl?: string
}>()

// 监听地址变化
onMounted(() => {
  if (props.address && props.utxoUrl) {
    handleGetUtxo(props.address)
  }
})
</script>

<style scoped>
.utxo-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.empty-state {
  margin: 40px 0;
}

.action-buttons {
  margin-top: 20px;
  text-align: center;
}
</style>

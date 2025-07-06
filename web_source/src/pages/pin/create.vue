<template>
  <div class="create-container">
    <div class="create-layout">
      <!-- 左侧：创建表单 -->
      <div class="form-section">
        <h2 class="section-title">创建 PIN</h2>
        <Pin ref="pinRef" />
      </div>

      <!-- 右侧：结果展示 -->
      <div class="result-section">
        <div class="result-wrapper">
          <!-- 创建结果 -->
          <div v-if="createResult" class="result-card">
            <h3 class="card-title">创建结果</h3>
            <el-descriptions :column="1" border>
              <el-descriptions-item>
                <template #label>
                  <div class="result-label">Commit Transaction</div>
                </template>
                <div class="tx-content">
                  <el-input v-model="createResult.commit_tx_raw" type="textarea" :rows="2" readonly />
                  <el-button type="primary" link @click="copyContent(createResult.commit_tx_raw)">
                    <el-icon><Document /></el-icon>
                    复制
                  </el-button>
                </div>
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>
                  <div class="result-label">Reveal Transaction</div>
                </template>
                <div class="tx-content">
                  <el-input v-model="createResult.reveal_tx_raw" type="textarea" :rows="2" readonly />
                  <el-button type="primary" link @click="copyContent(createResult.reveal_tx_raw)">
                    <el-icon><Document /></el-icon>
                    复制
                  </el-button>
                </div>
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>
                  <div class="result-label">Miner Fee</div>
                </template>
                <div class="fee-content">
                  {{ createResult.miner_fee / 100000000 }} BTC
                </div>
              </el-descriptions-item>
            </el-descriptions>
          </div>

          <!-- 广播结果 -->
          <div v-if="broadcastResult.commit || broadcastResult.reveal" class="result-card">
            <h3 class="card-title">广播结果</h3>
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
        </div>
      </div>
    </div>

    <!-- UTXO 抽屉 -->
    <el-drawer
      v-model="utxoDrawerVisible"
      title="选择 UTXO"
      direction="rtl"
      size="50%"
    >
      <GetUtxo 
        :address="currentAddress"
        :utxo-url="currentUtxoUrl"
        @select="handleSelectUtxo" 
      />
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import { ref, provide, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Document } from '@element-plus/icons-vue'
import Pin from '../../components/Pin.vue'
import GetUtxo from '../../components/GetUtxo.vue'
import type { Utxo } from '../../api/modules/inscribe/types'
import type { InscribeResponse } from '../../api/modules/inscribe/types'
import { broadcastTx } from '../../api/modules/inscribe'

const activeName = ref('pin')
const pinRef = ref()
const utxoDrawerVisible = ref(false)
const createResult = ref<InscribeResponse | null>(null)
const currentAddress = ref('')
const currentUtxoUrl = ref('')

// 广播结果
const broadcastResult = reactive({
  commit: '',
  reveal: ''
})

// 提供方法给子组件
provide('showUtxoDrawer', (address: string, utxoUrl: string) => {
  if (!address) {
    ElMessage.warning('请先选择地址')
    return
  }
  if (!utxoUrl) {
    ElMessage.warning('UTXO URL not configured')
    return
  }
  currentAddress.value = address
  currentUtxoUrl.value = utxoUrl
  utxoDrawerVisible.value = true
})

provide('setCreateResult', (result: InscribeResponse) => {
  createResult.value = result
})

const handleSelectUtxo = (utxos: Utxo[]) => {
  if (pinRef.value?.addUtxo) {
    utxos.forEach(utxo => {
      pinRef.value.addUtxo(utxo)
    })
  }
  utxoDrawerVisible.value = false
}

const handleBroadcast = async () => {
  if (!createResult.value) return
  
  try {
    // 广播 Commit 交易
    const commitResult = await broadcastTx(createResult.value.commit_tx_raw)
    ElMessage.success('Commit 交易广播成功')
    broadcastResult.commit = commitResult.txId.replace(/"/g, '')
    
    // 广播 Reveal 交易
    const revealResult = await broadcastTx(createResult.value.reveal_tx_raw)
    ElMessage.success('Reveal 交易广播成功')
    broadcastResult.reveal = revealResult.txId.replace(/"/g, '')
    
    // 不再清空创建结果
    // createResult.value = null
  } catch (error) {
    ElMessage.error('广播交易失败')
    console.error('Error:', error)
  }
}

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
</script>

<style scoped>
.create-container {
  padding: 24px;
  height: 100%;
}

.create-layout {
  display: flex;
  gap: 24px;
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
}

.form-section {
  flex: 1;
  min-width: 0;
  border-radius: 8px;
  padding: 24px;
}

.result-section {
  width: 500px;
  max-height: calc(100vh - 48px);
  overflow-y: auto;
  padding-right: 8px;
  display: flex;
  flex-direction: column;
}

.result-wrapper {
  display: flex;
  flex-direction: column;
  gap: 24px;
  min-height: 100%;
}

.result-section::-webkit-scrollbar {
  width: 6px;
}

.result-section::-webkit-scrollbar-thumb {
  background-color: var(--el-border-color);
  border-radius: 3px;
}

.result-section::-webkit-scrollbar-track {
  background-color: var(--el-bg-color);
}

.section-title {
  margin: 0 0 24px;
  font-size: 20px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.result-card {
  border-radius: 8px;
  padding: 24px;
  border: 1px solid var(--el-border-color-light);
  background-color: var(--el-bg-color);
  width: 100%;
  box-sizing: border-box;
}

.card-title {
  margin: 0 0 20px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.result-label {
  font-weight: 500;
  color: var(--el-text-color-regular);
  margin-bottom: 8px;
}

.tx-content {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.tx-content .el-input {
  flex: 1;
}

.tx-content .el-button {
  margin-top: 4px;
}

.fee-content {
  font-size: 16px;
  color: var(--el-text-color-primary);
  font-weight: 500;
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

<template>
  <div class="broadcast-container">
    <el-form :model="form" label-width="120px" :rules="rules" ref="formRef">
      <el-form-item label="Transaction" prop="txRaw">
        <el-input
          v-model="form.txRaw"
          type="textarea"
          :rows="4"
          placeholder="Enter transaction raw data"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleBroadcast">Broadcast</el-button>
      </el-form-item>
    </el-form>

    <div v-if="txId" class="result">
      <h3>Transaction ID:</h3>
      <el-input v-model="txId" readonly />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { broadcastTx } from '../api/modules/inscribe'

const formRef = ref<FormInstance>()
const txId = ref('')

const form = reactive({
  txRaw: ''
})

const rules = reactive<FormRules>({
  txRaw: [{ required: true, message: 'Please enter transaction raw data', trigger: 'blur' }]
})

const handleBroadcast = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const response = await broadcastTx(form.txRaw)
        txId.value = response.txId
        ElMessage.success('Transaction broadcasted successfully')
      } catch (error) {
        ElMessage.error('Failed to broadcast transaction')
        console.error('Error:', error)
      }
    }
  })
}
</script>

<style scoped>
.broadcast-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.result {
  margin-top: 20px;
}

.result h3 {
  margin-bottom: 10px;
}
</style>

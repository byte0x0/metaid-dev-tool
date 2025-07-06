<template>
  <div class="address-container">
    <div class="header">
      <el-button type="primary" @click="handleCreate">创建新地址</el-button>
    </div>

    <el-table :data="addressList" v-loading="loading" style="width: 100%">
      <el-table-column prop="chain.name" label="链名称" width="120" />
      <el-table-column prop="address" label="地址" min-width="300" />
      <el-table-column prop="private_key" label="私钥" min-width="300" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          {{ row.type === 'taproot' ? 'Taproot' : 'Segwit' }}
        </template>
      </el-table-column>
      <el-table-column prop="CreatedAt" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" width="120">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="创建新地址" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="链">
          <el-select v-model="form.chain_id" placeholder="请选择链">
            <el-option
              v-for="chain in chainList"
              :key="chain.ID"
              :label="chain.name"
              :value="chain.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="地址类型">
          <el-select v-model="form.type" placeholder="请选择地址类型">
            <el-option label="Taproot" :value="AddressType.Taproot" />
            <el-option label="Segwit" :value="AddressType.Segwit" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            创建
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAddressList, createAddress, deleteAddress } from '../../api/modules/address'
import { getChainList } from '../../api/modules/chain'
import type { AddressInfo } from '../../api/modules/address/types'
import type { ChainInfo } from '../../api/modules/chain/types'
import { AddressType } from '../../api/modules/address/types'
import { formatDate } from '../../utils/format'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const addressList = ref<AddressInfo[]>([])
const chainList = ref<ChainInfo[]>([])

const form = ref({
  chain_id: undefined as number | undefined,
  type: AddressType.Taproot
})

// 获取地址列表
const fetchAddressList = async () => {
  loading.value = true
  try {
    const res = await getAddressList({
      page: 1,
      pageSize: 10
    })
    addressList.value = res
  } catch (error) {
    console.error('获取地址列表失败:', error)
    ElMessage.error('获取地址列表失败')
  } finally {
    loading.value = false
  }
}

// 获取链列表
const fetchChainList = async () => {
  try {
    const res = await getChainList({})
    chainList.value = res
  } catch (error) {
    console.error('获取链列表失败:', error)
    ElMessage.error('获取链列表失败')
  }
}

// 创建地址
const handleCreate = () => {
  form.value = {
    chain_id: undefined,
    type: AddressType.Taproot
  }
  dialogVisible.value = true
}

// 提交创建
const handleSubmit = async () => {
  if (!form.value.chain_id) {
    ElMessage.warning('请选择链')
    return
  }

  submitting.value = true
  try {
    await createAddress({
      chain_id: form.value.chain_id,
      type: form.value.type
    })
    ElMessage.success('创建成功')
    dialogVisible.value = false
    fetchAddressList()
  } catch (error) {
    console.error('创建地址失败:', error)
    ElMessage.error('创建地址失败')
  } finally {
    submitting.value = false
  }
}

// 删除地址
const handleDelete = async (row: AddressInfo) => {
  try {
    await ElMessageBox.confirm('确定要删除该地址吗？', '提示', {
      type: 'warning'
    })
    await deleteAddress(row.ID)
    ElMessage.success('删除成功')
    fetchAddressList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除地址失败:', error)
      ElMessage.error('删除地址失败')
    }
  }
}

onMounted(() => {
  fetchAddressList()
  fetchChainList()
})
</script>

<style scoped>
.address-container {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>

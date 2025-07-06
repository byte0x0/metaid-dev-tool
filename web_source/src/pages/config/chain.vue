<template>
  <div class="chain-container">
    <div class="chain-header">
      <h2>链配置管理</h2>
      <el-button type="primary" @click="handleCreateChain">
        <el-icon><Plus /></el-icon>新建链
      </el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="chainList"
      style="width: 100%"
      border
    >
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column prop="chain_type" label="类型" min-width="100">
        <template #default="{ row }">
          {{ formatChainType(row.chain_type) }}
        </template>
      </el-table-column>
      <el-table-column prop="broadcast_url" label="广播URL" min-width="200" />
      <el-table-column prop="utxo_url" label="UTXO URL" min-width="200" />
      <el-table-column prop="CreatedAt" label="创建时间" min-width="180">
        <template #default="{ row }">
          {{ formatDate(row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column prop="UpdatedAt" label="更新时间" min-width="180">
        <template #default="{ row }">
          {{ formatDate(row.UpdatedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button type="primary" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              删除
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next"
      class="pagination"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新建链' : '编辑链'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入链名称" />
        </el-form-item>
        <el-form-item label="类型" prop="chainType">
          <el-select v-model="formData.chainType" placeholder="请选择链类型">
            <el-option label="主网" value="mainNet" />
            <el-option label="测试网" value="testNet" />
            <el-option label="回归测试网" value="regTest" />
          </el-select>
        </el-form-item>
        <el-form-item label="广播URL" prop="broadcastUrl">
          <el-input v-model="formData.broadcastUrl" placeholder="请输入广播URL" />
        </el-form-item>
        <el-form-item label="UTXO URL" prop="utxoUrl">
          <el-input v-model="formData.utxoUrl" placeholder="请输入UTXO URL（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getChainList,
  createChain,
  updateChain,
  deleteChain,
} from '../../api/modules/chain'
import type { ChainInfo, CreateChainParams, UpdateChainParams } from '../../api/modules/chain/types'

// 数据列表相关
const loading = ref(false)
const chainList = ref<ChainInfo[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 表单相关
const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const formData = ref<CreateChainParams>({
  name: '',
  chainType: 'mainNet',
  broadcastUrl: '',
  utxoUrl: '',
})
const currentChainId = ref('')

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入链名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  chainType: [
    { required: true, message: '请选择链类型', trigger: 'change' },
  ],
  broadcastUrl: [
    { required: true, message: '请输入广播URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' },
  ],
  utxoUrl: [
    { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' },
  ],
}

// 获取链列表
const fetchChainList = async () => {
  loading.value = true
  try {
    const res = await getChainList({})
    chainList.value = res
    total.value = res.length // 由于后端没有分页，暂时使用数组长度
  } catch (error) {
    ElMessage.error('获取链列表失败')
  } finally {
    loading.value = false
  }
}

// 处理分页
const handleSizeChange = (val: number) => {
  if (pageSize.value === val) return
  pageSize.value = val
  fetchChainList()
}

const handleCurrentChange = (val: number) => {
  if (currentPage.value === val) return
  currentPage.value = val
  fetchChainList()
}

// 处理创建
const handleCreateChain = () => {
  dialogType.value = 'create'
  formData.value = {
    name: '',
    chainType: 'mainNet',
    broadcastUrl: '',
    utxoUrl: '',
  }
  dialogVisible.value = true
}

// 处理编辑
const handleEdit = (row: ChainInfo) => {
  dialogType.value = 'edit'
  currentChainId.value = row.ID.toString()
  formData.value = {
    name: row.name,
    chainType: row.chain_type,
    broadcastUrl: row.broadcast_url,
    utxoUrl: row.utxo_url,
  }
  dialogVisible.value = true
}

// 处理删除
const handleDelete = async (row: ChainInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除链 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      },
    )
    await deleteChain(row.ID.toString())
    ElMessage.success('删除成功')
    fetchChainList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 处理提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'create') {
          await createChain(formData.value)
          ElMessage.success('创建成功')
        } else {
          await updateChain(currentChainId.value, formData.value)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        fetchChainList()
      } catch (error) {
        ElMessage.error(dialogType.value === 'create' ? '创建失败' : '更新失败')
      }
    }
  })
}

// 格式化日期
const formatDate = (date: string) => {
  if (!date) return '-'
  try {
    const d = new Date(date)
    if (isNaN(d.getTime())) return '-'
    return d.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
  } catch (e) {
    return '-'
  }
}

// 格式化链类型
const formatChainType = (type: string) => {
  const typeMap: Record<string, string> = {
    mainNet: '主网',
    testNet: '测试网',
    regTest: '回归测试网',
  }
  return typeMap[type] || type
}

// 初始化
onMounted(() => {
  fetchChainList()
})
</script>

<style scoped>
.chain-container {
  padding: 20px;
}

.chain-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-form-item__content) {
  width: 100%;
}
</style>

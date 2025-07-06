<template>
  <div style="width: 90%">
    <div style="text-align: left">
      <el-button type="primary" @click="create">Create New</el-button>
    </div>
    <el-table :data="tableData">
      <el-table-column prop="date" label="Chain" width="180" />
      <el-table-column prop="date" label="Net" width="180" />
      <el-table-column prop="name" label="Url" />
      <el-table-column prop="address" label="Example" />
      <el-table-column fixed="right" label="Operations" width="120">
        <template #default="scope">
          <el-button link type="primary" size="small" @click="deleteItem(scope.row)"> Delete </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
  <el-dialog v-model="dialogFormVisible" title="New address" width="800">
    <el-form :model="form">
      <el-form-item label="Chain" :label-width="formLabelWidth">
        <el-select v-model="form.chain" placeholder="Please select chain">
          <el-option label="BTC" value="BTC" />
          <el-option label="MVC" value="MVC" />
          <el-option label="FB" value="FB" />
        </el-select>
        <el-select v-model="form.net" placeholder="Please select net">
          <el-option label="regtest" value="regtest" />
          <el-option label="testnet" value="testnet" />
          <el-option label="mainnet" value="mainnet" />
        </el-select>
      </el-form-item>
      <el-form-item label="Url" :label-width="formLabelWidth">
        <el-input v-model="form.url" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Example" :label-width="formLabelWidth">
        <el-input v-model="form.example" autocomplete="off" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="dialogFormVisible = false"> Confirm </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'

const dialogFormVisible = ref(false)
const formLabelWidth = '140px'
const form = reactive({
  chain: '',
  net: '',
  url: '',
  example: '',
})
const tableData = [
  {
    date: '2016-05-03',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-02',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-04',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    date: '2016-05-01',
    name: 'Tom',
    address: 'No. 189, Grove St, Los Angeles',
  },
]
const create = () => {
  dialogFormVisible.value = true
}
const deleteItem = (v: any) => {}
</script>

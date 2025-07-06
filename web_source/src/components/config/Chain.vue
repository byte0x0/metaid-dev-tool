<template>
  <el-form label-width="auto" :model="formData">
    <el-form-item label="RegTest">
      <el-input v-model="formData.regtest" autosize type="textarea" />
    </el-form-item>
    <el-form-item label="TestNet">
      <el-input v-model="formData.testnet" autosize type="textarea" />
    </el-form-item>
    <el-form-item label="MainNet">
      <el-input v-model="formData.mainnet" autosize type="textarea" />
    </el-form-item>
    <el-button type="primary" @click="save">Save</el-button>
  </el-form>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'

const props = defineProps({
  chain: {
    type: String,
    default: '',
  },
})
const formData = reactive({
  chain: props.chain,
  regtest: '',
  testnet: '',
  mainnet: '',
})
const save = () => {
  console.log(formData)
}
</script>

<script lang="ts" setup>
import { Document, Menu as IconMenu, Location, Setting } from '@element-plus/icons-vue'

// const isCollapse = ref(true)
function handleOpen(key: string, keyPath: string[]) {
  // eslint-disable-next-line no-console
  console.log(key, keyPath)
}
function handleClose(key: string, keyPath: string[]) {
  // eslint-disable-next-line no-console
  console.log(key, keyPath)
}
</script>

<template>
  <el-menu router default-active="1" class="el-menu-vertical-demo" @open="handleOpen" @close="handleClose">
    <el-menu-item index="/pin/create">
      <el-icon>
        <IconMenu />
      </el-icon>
      <template #title> Create PIN </template>
    </el-menu-item>

    <el-menu-item index="/address">
      <el-icon>
        <Location />
      </el-icon>
      <template #title> Address </template>
    </el-menu-item>

    <el-menu-item index="/config/chain">
      <el-icon>
        <Setting />
      </el-icon>
      <template #title> Chain </template>
    </el-menu-item>

    
  </el-menu>
</template>

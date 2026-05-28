<template>
  <el-table :data="userStore.users" height="80%">
    <el-table-column label="用户名" prop="user_name" />

    <el-table-column label="昵称" prop="nickname" />

    <el-table-column label="是否为管理员">
      <template #default="scope">{{ scope.row.is_admin ? "是" : "否" }}</template>
    </el-table-column>

    <el-table-column label="注册时间">
      <template #default="scope">{{ displayTimestamp(scope.row.created_at) }}</template>
    </el-table-column>

    <el-table-column label="上次登录时间">
      <template #default="scope">{{ displayTimestamp(scope.row.last_login) }}</template>
    </el-table-column>
  </el-table>

  <el-pagination
    layout="prev,pager,next,->,total"
    :total="userStore.count"
    background
    @current-change="listUser"
  />
</template>

<script lang="ts" setup>
import { onMounted } from "vue"
import { useUserStore } from "@/pinia/user.ts"
import { displayTimestamp } from "@/ts/util.ts"

const userStore = useUserStore()

onMounted(() => {
    listUser()
})

function listUser(pageNum: number = 1): void {
    userStore.list(10, pageNum)
}
</script>

<style lang="less" scoped>
.el-pagination {
	height: 20%;
}
</style>

<template>
	<el-table :data="users" height="80%">
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

	<el-pagination layout="prev,pager,next,->,total" :total="amount" background @currentChange="listUser" />
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue"
import { User } from "@/axios/ts/user.go.ts"
import { useUserStore } from "@/pinia/user.ts"
import { displayTimestamp } from "@/ts/util.ts"

let userStore = useUserStore()

let amount = ref<number>(0)
let users = ref<Array<User>>(new Array<User>())

onMounted(() => {
	listUser()
})

function listUser(pageNum: number = 1): void {
	userStore.list(10, pageNum, (a: number, u: Array<User>) => {
		amount.value = a
		users.value = u
	})
}
</script>

<style lang="less" scoped>
.el-pagination {
	height: 20%;
}
</style>

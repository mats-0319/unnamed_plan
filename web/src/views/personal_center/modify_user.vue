<template>
  <el-form v-model="modifyUserReq" class="modify-user" label-width="20%">
    <el-form-item label="用户名">{{ userStore.user.user_name }}</el-form-item>

    <el-form-item label="昵称"><el-input v-model="modifyUserReq.nickname" /></el-form-item>

    <el-form-item label="密码">
      <el-input v-model="modifyUserReq.password" show-password />
    </el-form-item>

    <el-form-item>
      <outlined-button :details="tips_ModifyUser" :disabled="!canModifyFlag" @click="modifyUser()">
        修改个人信息
      </outlined-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
import { ModifyUserReq } from "@/axios/ts/user.go.ts"
import { onMounted, ref, watch } from "vue"
import { useUserStore } from "@/pinia/user.ts"
import OutlinedButton from "@/components/outlined_button.vue"
import { tips_ModifyUser } from "@/ts/data.ts"

let userStore = useUserStore()

let modifyUserReq = ref<ModifyUserReq>(new ModifyUserReq())
let canModifyFlag = ref<boolean>(false)

onMounted(() => {
    modifyUserReq.value.nickname = userStore.user.nickname
})

function modifyUser(): void {
    userStore.modify(
        modifyUserReq.value.nickname,
        modifyUserReq.value.password,
    )
}

watch(
    modifyUserReq,
    (newValue, _) => {
        canModifyFlag.value =
            newValue.nickname != userStore.user.nickname || newValue.password.length > 0
    },
    { deep: true },
)
</script>

<style lang="less">
.modify-user {
	.el-input {
		width: 50%;
	}
}
</style>

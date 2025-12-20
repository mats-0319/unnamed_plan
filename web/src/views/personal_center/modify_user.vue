<template>
  <el-form class="modify-user" v-model="modifyUserReq" label-width="20%">
    <el-form-item label="用户名">
      {{ userStore.user.name }}
    </el-form-item>

    <el-form-item label="昵称">
      <el-input v-model="modifyUserReq.nickname"/>
    </el-form-item>

    <el-form-item label="密码">
      <el-input v-model="modifyUserReq.password" show-password/>
    </el-form-item>

    <el-form-item label="是否修改TOTP密钥">
      <el-switch v-model="modifyUserReq.modify_tk_flag"/>
      &emsp;{{ modifyUserReq.modify_tk_flag ? "修改" : "不修改" }}
    </el-form-item>

    <el-form-item label="TOTP密钥">
      <el-input v-model="modifyUserReq.totp_key"/>
    </el-form-item>

    <el-form-item>
      <up-button
          details="密码为空表示不修改<br/>
          修改TOTP密钥且新值为空，表示关闭TOTP功能<br/><br/>
          昵称、密码、TOTP均无修改时，不可提交"
          :disabled="!canModifyFlag"
          @click="modifyUser()"
      >
        修改个人信息
      </up-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
import {ModifyUserReq} from "@/axios/ts/user.go.ts";
import {onMounted, ref, watch} from "vue";
import {useUserStore} from "@/pinia/user.ts";
import UpButton from "@/components/up_button.vue";

let userStore = useUserStore();

let modifyUserReq = ref<ModifyUserReq>(new ModifyUserReq());
let canModifyFlag = ref<boolean>(false);

onMounted(() => {
  modifyUserReq.value.nickname = userStore.user.nickname
})

function modifyUser(): void {
  userStore.modify(modifyUserReq.value.nickname,
      modifyUserReq.value.password,
      modifyUserReq.value.modify_tk_flag,
      modifyUserReq.value.totp_key)
}

watch(modifyUserReq, (newValue, oldValue) => {
  canModifyFlag.value = newValue.nickname != userStore.user.nickname ||
      newValue.password.length > 0 ||
      newValue.modify_tk_flag
}, {deep: true})
</script>

<style lang="less">
.modify-user {
  .el-input {
    width: 50%;
  }
}
</style>

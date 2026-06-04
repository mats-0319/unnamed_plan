<template>
  <el-form v-model="setMFAStatusReq" class="set-mfa-status" labelWidth="20%" size="large">
    <el-form-item label="是否启用MFA">
      <el-switch v-model="setMFAStatusReq.enable_mfa" />&emsp;
      {{ setMFAStatusReq.enable_mfa ? "是" : "否" }}
    </el-form-item>

    <template v-if="setMFAStatusReq.enable_mfa">
      <el-form-item label="TOTP密钥">
        <span v-if="setMFAStatusReq.apply_new_key_flag">{{ totpKey }}</span>
        <template v-else>
          <span class="flex">
            <b v-if="userStore.user.has_totp_key"><i>[继续使用历史密钥]&emsp;</i></b>
            <elevatedButton :onClick="applyNewTOTPKey">申请新的密钥</elevatedButton>
          </span>
        </template>
      </el-form-item>

      <el-form-item label="TOTP Code">
        <el-input-otp v-model="setMFAStatusReq.totp_code" type="filled" />
      </el-form-item>
    </template>

    <el-form-item>
      <elevated-button :onClick="setMFAStatus">设置MFA状态</elevated-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
import { useUserStore } from "@/pinia/user.ts"
import { SetMFAStatusReq } from "@/axios/ts/user.go.ts"
import { onMounted, ref } from "vue"
import ElevatedButton from "@/components/elevated_button.vue"

const userStore = useUserStore()

const totpKey = ref<string>("")
const setMFAStatusReq = ref<SetMFAStatusReq>(new SetMFAStatusReq())

onMounted(() => {
    setMFAStatusReq.value = new SetMFAStatusReq()
    setMFAStatusReq.value.enable_mfa = userStore.user.enable_mfa
})

async function applyNewTOTPKey() {
    const res = await userStore.applyTOTPKey()

    totpKey.value = res.totp_key
    setMFAStatusReq.value.apply_new_key_flag = true
}

async function setMFAStatus() {
    await userStore.setMFAStatus(
        setMFAStatusReq.value.enable_mfa,
        setMFAStatusReq.value.apply_new_key_flag,
        setMFAStatusReq.value.totp_code,
    )
}
</script>

<style lang="less" scoped>
.set-mfa-status {
  .flex {
    display: flex;
    font-size: 1.2rem;
  }

  .tips {
    opacity: 0.6;
  }
}
</style>

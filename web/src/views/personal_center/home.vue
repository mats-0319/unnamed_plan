<template>
  <div class="personal-center color-bg-0">
    <div class="p-left">
      <div
          v-if="userStore.user.is_admin"
          class="pl-item color-bg-1 center-hv"
          @click="routerLink('pListUser')"
      >
        查看用户
      </div>

      <div class="pl-item color-bg-1 center-hv" @click="routerLink('pModifyUser')">
        修改个人信息
      </div>

      <div class="pl-item  color-bg-1 center-hv" @click="routerLink('pCreateNote')">
        编辑小纸条
      </div>
    </div>

    <el-divider direction="vertical"></el-divider>

    <div class="p-right">
      <router-view/>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {routerLink} from "@/ts/util.ts";
import {useUserStore} from "@/pinia/user.ts";

let userStore = useUserStore()
</script>

<style lang="less" scoped>
.personal-center {
  display: flex;
  height: calc(100vh - 6.25rem);

  .p-left {
    width: calc(20vw - 1px);
    margin: 6rem 2vw;
    font-size: 1.4rem;

    .pl-item {
      width: 10vw;
      margin: 1.5rem auto;
      padding: 0.5rem 0;
      border-radius: 1rem;
    }
    .pl-item:hover {
      cursor: pointer;
      text-decoration-line: underline;
    }
  }

  .p-right {
    width: 60vw;
    margin: 8rem 8vw;
  }

  .el-divider--vertical {
    width: 0;
    height: 100%;
    margin: 0;
    border-left: 1px solid darkgray;
  }
}
</style>

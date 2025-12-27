<template>
  <div class="personal-center color-bg-0">
    <div class="p-left">
      <elevated_button
          class="pl-item"
          v-if="userStore.user.is_admin"
          @click="routerLink('pListUser')"
      >
        查看用户
      </elevated_button>

      <elevated_button class="pl-item" @click="routerLink('pModifyUser')">
        修改个人信息
      </elevated_button>

      <elevated_button class="pl-item" @click="routerLink('pNote')">
        我的小纸条
      </elevated_button>
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
import Elevated_button from "@/components/elevated_button.vue";

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
      margin: 2rem auto;
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

<template>
  <div class="note color-bg-0">
    note
  </div>

  <bottom/>
</template>

<script lang="ts" setup>

import Bottom from "@/components/bottom.vue";
</script>

<style lang="less" scoped>
.note {
  height: calc(100vh - 6.25rem - 6.25rem);
}
</style>

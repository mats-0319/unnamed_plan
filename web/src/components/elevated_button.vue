<template>
  <div class="elevated-button" :style="{ width: props.w + 'rem', height: props.h + 'rem' }">
    <lock-button
      :style="{ 'background-color': props.bg }"
      :disabled="props.disabled"
      :on-click="props.onClick"
    >
      <slot />
    </lock-button>
  </div>
</template>

<script lang="ts" setup>
import LockButton from "@/components/lock_button.vue"

const props = defineProps({
    w: {
        type: String,
        default: "fit-content",
    },
    h: {
        type: String,
        default: "fit-content",
    },
    bg: {
        type: String,
        default: "rgb(240, 239, 226)",
    },
    disabled: {
        type: Boolean,
        default: false,
    },
    onClick: {
        type: Function,
        required: true,
    },
})
</script>

<style lang="less" scoped>
.elevated-button{
  border-radius: 1rem;

  button {
    width: 100%;
    height: 100%;
    border: none;
    border-radius: 1rem;
  }
  button:hover,
  button.loading {
    border: 1px solid darkgray;
    opacity: 0.8;
  }
}
</style>

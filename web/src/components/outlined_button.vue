<template>
  <div class="outlined-button" :style="{ width: props.w + 'rem', height: props.h + 'rem' }">
    <lock-button :disabled="props.disabled" :on-click="props.onClick"><slot /></lock-button>
  </div>
</template>

<script setup lang="ts">
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

<style lang="less">
.outlined-button {
  button {
    width: 100%;
    height: 100%;
    background-color: transparent;
    border: 1px solid black
  }
  button:hover {
    opacity: 0.8;
  }
  button:disabled {
    border-color: darkgray;
  }
}
</style>

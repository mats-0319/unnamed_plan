<template>
  <div class="up-button">
    <div class="u-box center-hv" :class="{ 'u-disabled' : props.disabled}">
      &emsp;<span class="ub-text"><slot></slot></span>&emsp;
    </div>

    <div v-if="props.details.length > 0" class="u-details">
      <el-tooltip :content="props.details" placement="top" raw-content>
        &nbsp;?&nbsp;
      </el-tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
let props = defineProps({
  details: {
    type: String,
    default: "",
  },
  disabled: {
    type: Boolean,
    default: false,
  }
})
</script>

<style lang="less">
.up-button {
  display: flex;

  .u-box {
    border: 1px solid darkgrey;
  }

  .u-box:hover {
    cursor: pointer;

    .ub-text {
      text-decoration-line: underline;
    }
  }

  .u-disabled {
    pointer-events: none;
    opacity: 0.6;
  }

  .u-details {
    margin-left: 1rem;

    .el-tooltip__trigger {
      border-radius: 50%;
      border: 1px solid darkgray;
    }
  }
}
</style>
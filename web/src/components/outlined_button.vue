<template>
	<div class="outlined-button">
		<div class="ob-box center-hv" :class="{ 'ob-disabled': props.disabled }">
			&emsp;<span class="obb-text"><slot></slot></span>&emsp;
		</div>

		<div v-if="props.details.length > 0" class="ob-details">
			<el-tooltip :content="props.details" placement="top" raw-content> &nbsp;?&nbsp; </el-tooltip>
		</div>
	</div>
</template>

<script setup lang="ts">
let props = defineProps({
	details: {
		type: String,
		default: ""
	},
	disabled: {
		type: Boolean,
		default: false
	}
})
</script>

<style lang="less">
.outlined-button {
	display: flex;

	.ob-box {
		border: 1px solid darkgrey;
	}

	.ob-box:hover {
		cursor: pointer;

		.obb-text {
			text-decoration-line: underline;
		}
	}

	.ob-disabled {
		pointer-events: none;
		opacity: 0.6;
	}

	.ob-details {
		margin-left: 1rem;

		.el-tooltip__trigger {
			border-radius: 50%;
			border: 1px solid darkgray;
		}
	}
}
</style>

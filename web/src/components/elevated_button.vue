<template>
	<div
		class="elevated-button color-bg-1 center-hv"
		:style="{ width: props.width + 'rem' }"
		:class="{ 'eb-loading': props.loading }"
	>
		<span v-show="!loading"><slot></slot></span>
		<span v-show="loading">loading...</span>
	</div>
</template>

<script lang="ts" setup>
let props = defineProps({
	width: {
		type: Number,
		default: 12
	},
	loading: {
		type: Boolean,
		default: false
	}
})
</script>

<style lang="less" scoped>
.elevated-button {
	border-radius: 1rem;
	padding: 0.5rem 0;
}
.elevated-button:hover {
	cursor: pointer;
	text-decoration-line: underline;
}

.eb-loading {
	pointer-events: none;
	opacity: 0.6;
}
</style>

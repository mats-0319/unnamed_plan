<template>
	<div class="note color-bg-0 center-hv">
		<stacked_cards class="n-card color-bg-1">
			<div class="nc-content">
				<p class="ncc-title">
					小纸条<span class="ncct-right">{{ currentNum }} / {{ amount }}</span>
				</p>
				<el-form class="ncc-form" v-model="currentNote" label-width="20%">
					<el-form-item label="ID">{{ currentNote.note_id }}</el-form-item>
					<el-form-item label="作者">{{ currentNote.is_anonymous ? "-" : currentNote.writer }}</el-form-item>
					<el-form-item label="主题">{{ currentNote.title }}</el-form-item>
					<el-form-item label="内容">{{ currentNote.content }}</el-form-item>
					<el-form-item label="写作时间">{{ displayTimestamp(currentNote.created_at) }}</el-form-item>
					<el-form-item label="修改时间">{{ displayTimestamp(currentNote.updated_at) }}</el-form-item>
				</el-form>
			</div>
		</stacked_cards>

		<div class="n-options">
			<elevated_button :loading="nextLoadingFlag" @click="nextNote">下一张</elevated_button>
		</div>
	</div>

	<bottom />
</template>

<script lang="ts" setup>
import Bottom from "@/components/bottom.vue"
import Stacked_cards from "@/components/stacked_cards.vue"
import { onMounted, ref } from "vue"
import { Note } from "@/axios/ts/note.go.ts"
import { useNoteStore } from "@/pinia/note.ts"
import { displayTimestamp } from "@/ts/util.ts"
import Elevated_button from "@/components/elevated_button.vue"

let noteStore = useNoteStore()

let amount = ref<number>(0)
let nextPage = ref<number>(1) // page num
let notes = ref<Array<Note>>(new Array<Note>())
let currentNote = ref<Note>(new Note())
let currentNum = ref<number>(0)

let nextLoadingFlag = ref<boolean>(false)

onMounted(() => {
	listNote()
})

function nextNote(): void {
	nextLoadingFlag.value = true

	let currentIndex = -1
	for (let i = 0; i < notes.value.length; i++) {
		if (currentNote.value.id == notes.value[i].id) {
			currentIndex = i
			break
		}
	}

	// 如果当前已经是列表的最后一个小纸条了，则获取下一页
	if (currentIndex >= 0 && currentIndex == notes.value.length - 1) {
		if (currentNum.value >= amount.value) {
			listNote()
			currentNum.value = 0
		} else {
			listNote(nextPage.value)
		}
		return
	}

	// 这里可以兼容`currentIndex=-1`的情况，不需要单独处理
	setCurrentNote(notes.value[currentIndex + 1])
}

function setCurrentNote(note: Note): void {
	currentNote.value = note
	currentNum.value++

	nextLoadingFlag.value = false
}

function listNote(pageNum: number = 1): void {
	noteStore.list(10, pageNum, false, (a: number, n: Array<Note>) => {
		amount.value = a
		notes.value = n
		nextPage.value = pageNum + 1

		setCurrentNote(n.length > 0 ? n[0] : new Note())
	})
}
</script>

<style lang="less" scoped>
.note {
	height: calc(100vh - 6.25rem - 6.25rem);

	.n-card {
		@media (min-width: 1280px) {
			width: 80rem;
		}
		width: 40rem;
		height: calc(100vh - 6.25rem - 6.25rem - 10rem);

		position: relative;
		z-index: 3;

		.nc-content {
			padding: 2rem 4rem;

			.ncc-title {
				font-style: italic;
				font-weight: 600;
				height: 2rem;

				.ncct-right {
					float: right;
				}
			}

			.ncc-form {
				max-height: calc(100vh - 6.25rem - 6.25rem - 10rem - 8rem - 4rem);
				overflow-y: auto;
				white-space: pre-wrap;

				&::-webkit-scrollbar {
					display: none;
				}
			}
		}
	}

	.n-options {
		@media (min-width: 1280px) {
			bottom: 8rem;
		}

		position: absolute;
		right: 1rem;
		bottom: 3rem;
	}
}
</style>

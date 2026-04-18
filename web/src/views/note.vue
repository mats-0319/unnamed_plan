<template>
  <div class="note color-bg-0 center-hv">
    <stacked-cards class="n-card color-bg-1">
      <div class="nc-content">
        <p class="ncc-title">
          小纸条<span class="ncct-right">{{ currentNum }} / {{ count }}</span>
        </p>

        <el-form v-model="currentNote" class="ncc-form" label-width="20%">
          <el-form-item label="ID">{{ currentNote.note_id }}</el-form-item>
          <el-form-item label="作者">{{ currentNote.is_anonymous ? "-" : currentNote.writer }}</el-form-item>
          <el-form-item label="主题">{{ currentNote.title }}</el-form-item>
          <el-form-item label="内容">{{ currentNote.content }}</el-form-item>
          <el-form-item label="写作时间">{{ displayTimestamp(currentNote.created_at) }}</el-form-item>
          <el-form-item label="修改时间">{{ displayTimestamp(currentNote.updated_at) }}</el-form-item>
        </el-form>
      </div>
    </stacked-cards>

    <elevated-button class="n-options" :loading="nextLoadingFlag" @click="nextNote">下一张</elevated-button>
  </div>

  <bottom />
</template>

<script lang="ts" setup>
import Bottom from "@/components/bottom.vue"
import StackedCards from "@/components/stacked_cards.vue"
import { onMounted, ref } from "vue"
import { Note } from "@/axios/ts/note.go.ts"
import { useNoteStore } from "@/pinia/note.ts"
import { displayTimestamp } from "@/ts/util.ts"
import ElevatedButton from "@/components/elevated_button.vue"

let noteStore = useNoteStore()

let count = ref<number>(0)
let notes = ref<Array<Note>>(new Array<Note>())

let nextPage = ref<number>(1) // page num
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
        if (currentNote.value.note_id == notes.value[i].note_id) {
            currentIndex = i
            break
        }
    }

    // 下一张
    if (currentIndex >= 0 && currentIndex == notes.value.length - 1) { // 如果当前已经是列表的最后一个小纸条了
        if (currentNum.value >= count.value) { // 没有下一页了，从头开始
            listNote()
            currentNum.value = 0
        } else { // 还有下一页
            listNote(nextPage.value)
        }
    } else {
        setCurrentNote(notes.value[currentIndex + 1]) // 这里可以兼容`currentIndex=-1`的情况，不需要单独处理
    }
}

function setCurrentNote(note: Note): void {
    currentNote.value = note
    currentNum.value++

    nextLoadingFlag.value = false
}

function listNote(pageNum: number = 1): void {
    noteStore.list(false, 10, pageNum, (c: number, n: Array<Note>) => {
        count.value = c
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

		width: 10rem;
		position: absolute;
		right: 1rem;
		bottom: 3rem;
	}
}
</style>

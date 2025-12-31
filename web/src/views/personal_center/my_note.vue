<template>
	<div class="pnote-tools">
		<elevated_button :width="8" @click="beforeCreate()">写小纸条</elevated_button>
	</div>

	<el-table :data="notes" height="60%">
		<el-table-column type="expand">
			<template #default="scope">
				<el-descriptions class="pnote-table-description" title="小纸条详情" border column="1" label-width="15%">
					<el-descriptions-item label="小纸条ID">{{ scope.row.note_id }}</el-descriptions-item>
					<el-descriptions-item label="作者昵称">{{ scope.row.writer }}</el-descriptions-item>
					<el-descriptions-item label="主题">{{ scope.row.title }}</el-descriptions-item>
					<el-descriptions-item label="内容">{{ scope.row.content }}</el-descriptions-item>
					<el-descriptions-item label="写作时间">
						{{ displayTimestamp(scope.row.created_at) }}
					</el-descriptions-item>
					<el-descriptions-item label="修改时间">
						{{ displayTimestamp(scope.row.updated_at) }}
					</el-descriptions-item>
				</el-descriptions>
			</template>
		</el-table-column>

		<el-table-column label="作者昵称" prop="writer" :min-width="2" />

		<el-table-column label="主题" prop="title" :min-width="5" show-overflow-tooltip />

		<el-table-column label="操作" :min-width="3">
			<template #default="scope">
				<div class="buttons-box">
					<outlined_button class="button-item" @click="beforeModify(scope.row)">编辑</outlined_button>
					<outlined_button @click="beforeDelete(scope.row)">删除</outlined_button>
				</div>
			</template>
		</el-table-column>
	</el-table>

	<el-pagination layout="prev,pager,next,->,total" :total="amount" background @currentChange="listNote" />

	<el-dialog v-model="showCreateDialog" title="写小纸条">
		<el-form v-model="createNoteReq" label-width="20%">
			<el-form-item label="主题">
				<el-input v-model="createNoteReq.title" />
			</el-form-item>

			<el-form-item label="内容">
				<el-input v-model="createNoteReq.content" type="textarea" :rows="4" />
			</el-form-item>

			<el-form-item label="是否匿名发布">
				<el-switch v-model="createNoteReq.is_anonymous" />
				&emsp;{{ createNoteReq.is_anonymous ? "匿名" : "不匿名" }}
			</el-form-item>
		</el-form>

		<template #footer>
			<el-button @click="showCreateDialog = false">取消</el-button>
			<el-button type="primary" :disabled="!canCreateFlag" @click="create()">提交</el-button>
		</template>
	</el-dialog>

	<el-dialog v-model="showModifyDialog" title="编辑小纸条">
		<el-form v-model="modifyNoteReq" label-width="20%">
			<el-form-item label="小纸条ID">{{ originData.note_id }}</el-form-item>

			<el-form-item label="主题">
				<el-input v-model="modifyNoteReq.title" />
			</el-form-item>

			<el-form-item label="内容">
				<el-input v-model="modifyNoteReq.content" type="textarea" :rows="4" />
			</el-form-item>

			<el-form-item label="是否匿名发布">
				<el-switch v-model="modifyNoteReq.is_anonymous" />
				&emsp;{{ createNoteReq.is_anonymous ? "匿名" : "不匿名" }}
			</el-form-item>
		</el-form>

		<template #footer>
			<el-button @click="showModifyDialog = false">取消</el-button>
			<el-button type="primary" :disabled="!canModifyFlag" @click="modify()">提交</el-button>
		</template>
	</el-dialog>

	<el-dialog v-model="showDeleteDialog" title="删除小纸条">
		<el-form label-width="20%">
			<el-form-item><b>是否确认删除下方小纸条?</b></el-form-item>
			<el-form-item label="小纸条ID">{{ originData.note_id }}</el-form-item>
			<el-form-item label="主题">{{ originData.title }}</el-form-item>
			<el-form-item label="内容">{{ originData.content }}</el-form-item>
		</el-form>

		<template #footer>
			<el-button @click="showDeleteDialog = false">取消</el-button>
			<el-button type="primary" @click="del()">删除</el-button>
		</template>
	</el-dialog>
</template>

<script lang="ts" setup>
import Elevated_button from "@/components/elevated_button.vue"
import { onMounted, ref, watch } from "vue"
import { CreateNoteReq, DeleteNoteReq, ModifyNoteReq, Note } from "@/axios/ts/note.go.ts"
import { deepCopy, displayTimestamp } from "@/ts/util.ts"
import { useNoteStore } from "@/pinia/note.ts"
import Outlined_button from "@/components/outlined_button.vue"

let noteStore = useNoteStore()

let amount = ref<number>(0)
let notes = ref<Array<Note>>(new Array<Note>())

let showCreateDialog = ref<boolean>(false)
let canCreateFlag = ref<boolean>(false)
let createNoteReq = ref<CreateNoteReq>(new CreateNoteReq())

let showModifyDialog = ref<boolean>(false)
let canModifyFlag = ref<boolean>(false)
let originData = ref<Note>(new Note()) // use for modify and delete
let modifyNoteReq = ref<ModifyNoteReq>(new ModifyNoteReq())

let showDeleteDialog = ref<boolean>(false)
let deleteNoteReq = ref<DeleteNoteReq>(new DeleteNoteReq())

onMounted(() => {
	listNote()
})

function beforeCreate(): void {
	createNoteReq.value = new CreateNoteReq()
	showCreateDialog.value = true
}

function create(): void {
	noteStore.create(createNoteReq.value.is_anonymous, createNoteReq.value.title, createNoteReq.value.content, () => {
		listNote()
		showCreateDialog.value = false
	})
}

function beforeModify(note: Note): void {
	originData.value = deepCopy(note)
	modifyNoteReq.value = {
		id: note.id,
		is_anonymous: note.is_anonymous,
		title: note.title,
		content: note.content
	}

	showModifyDialog.value = true
}

function modify(): void {
	noteStore.modify(
		modifyNoteReq.value.id,
		modifyNoteReq.value.is_anonymous,
		modifyNoteReq.value.title,
		modifyNoteReq.value.content,
		() => {
			listNote()
			showModifyDialog.value = false
		}
	)
}

function beforeDelete(note: Note): void {
	originData.value = deepCopy(note)
	deleteNoteReq.value = {
		id: note.id
	}

	showDeleteDialog.value = true
}

function del(): void {
	noteStore.del(deleteNoteReq.value.id, () => {
		listNote()
		showDeleteDialog.value = false
	})
}

function listNote(pageNum: number = 1): void {
	noteStore.list(10, pageNum, true, (a: number, n: Array<Note>) => {
		amount.value = a
		notes.value = n
	})
}

watch(
	createNoteReq,
	(newValue, _) => {
		canCreateFlag.value = newValue.content.length > 0
	},
	{ deep: true }
)

watch(
	modifyNoteReq,
	(newValue, _) => {
		canModifyFlag.value =
			newValue.is_anonymous != originData.value.is_anonymous ||
			newValue.title != originData.value.title ||
			newValue.content != originData.value.content
	},
	{ deep: true }
)
</script>

<style lang="less" scoped>
.pnote-tools {
	height: 4rem;
}

.el-pagination {
	height: calc(40% - 4rem);
}

.pnote-table-description {
	padding: 1rem 4rem;
	white-space: pre-wrap;
}

.buttons-box {
	display: flex;

	.button-item {
		margin-right: 0.5rem;
	}
}
</style>

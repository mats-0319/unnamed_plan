import { defineStore } from "pinia"
import { noteAxios } from "@/axios/ts/note.http.ts"
import { CreateNoteRes, DeleteNoteRes, ListNoteRes, ModifyNoteRes, Note } from "@/axios/ts/note.go.ts"
import { log } from "@/ts/log.ts"

export let useNoteStore = defineStore("note", () => {
	function create(isAnonymous: boolean, title: string, content: string, cb: () => void): void {
		noteAxios.createNote(isAnonymous, title, content).then(({}: { data: CreateNoteRes }) => {
			cb()

			log.success("create note")
		})
	}

	function list(
		pageSize: number,
		pageNum: number,
		listMyFlag: boolean,
		cb: (amount: number, notes: Array<Note>) => void
	): void {
		noteAxios.listNote({ size: pageSize, num: pageNum }, listMyFlag).then(({ data }: { data: ListNoteRes }) => {
			cb(data.amount, data.notes)

			log.success("list note")
		})
	}

	function modify(id: number, isAnonymous: boolean, title: string, content: string, cb: () => void): void {
		noteAxios.modifyNote(id, isAnonymous, title, content).then(({}: { data: ModifyNoteRes }) => {
			cb()

			log.success("modify note")
		})
	}

	function del(id: number, cb: () => void): void {
		noteAxios.deleteNote(id).then(({}: { data: DeleteNoteRes }) => {
			cb()

			log.success("delete note")
		})
	}

	return { create, list, modify, del }
})

import { defineStore } from "pinia"
import { noteAxios } from "@/axios/ts/note.http.ts"
import { CreateNoteRes, DeleteNoteRes, ListNoteRes, ModifyNoteRes, Note } from "@/axios/ts/note.go.ts"
import { log } from "@/ts/log.ts"

export const useNoteStore = defineStore("note", () => {
    function create(isAnonymous: boolean, title: string, content: string, cb: () => void): void {
        noteAxios.createNote(isAnonymous, title, content).then(({}: { data: CreateNoteRes }) => {
            cb()

            log.success("create note")
        })
    }

    function list(
        onlyOperator: boolean,
        pageSize: number,
        pageNum: number,
        cb: (amount: number, notes: Array<Note>) => void,
    ): void {
        noteAxios.listNote(onlyOperator, { size: pageSize, num: pageNum }).then(({ data }: { data: ListNoteRes }) => {
            cb(data.amount, data.notes)

            log.success("list note")
        })
    }

    function modify(noteID: string, isAnonymous: boolean, title: string, content: string, cb: () => void): void {
        noteAxios.modifyNote(noteID, isAnonymous, title, content).then(({}: { data: ModifyNoteRes }) => {
            cb()

            log.success("modify note")
        })
    }

    function del(noteID: string, cb: () => void): void {
        noteAxios.deleteNote(noteID).then(({}: { data: DeleteNoteRes }) => {
            cb()

            log.success("delete note")
        })
    }

    return { create, list, modify, del }
})

package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyNote(ctx *mhttp.Context) {
	req := &api.ModifyNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(req.NoteID) < 1 {
		e := ErrInvalidParams().WithParam("note id", req.NoteID)
		ctx.ResData = e
		mlog.Error(e.String())
		return
	}

	note, e := dal.GetNote(req.NoteID)
	if e != nil {
		ctx.ResData = e
		return
	}

	// no changes
	if req.IsAnonymous == note.IsAnonymous && req.Title == note.Title && req.Content == note.Content {
		e := ErrNoChanges().WithParam("operator", ctx.UserName)
		ctx.ResData = e
		mlog.Error(e.String())
		return
	}

	if ctx.UserName != note.Writer {
		e := ErrPermissionDenied().WithParam("need writer but get", ctx.UserName)
		ctx.ResData = e
		mlog.Error(e.String())
		return
	}

	note.IsAnonymous = req.IsAnonymous
	note.Title = req.Title
	note.Content = req.Content

	if e := dal.UpdateNote(note); e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ModifyNoteRes{}
}

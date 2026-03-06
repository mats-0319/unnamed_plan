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

	note, err := dal.GetNote(req.NoteID)
	if err != nil {
		ctx.ResData = err
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
		e := ErrNeedOwner().WithParam("operator", ctx.UserName).WithParam("owner", note.Writer)
		ctx.ResData = e
		mlog.Error(e.String())
		return
	}

	note.IsAnonymous = req.IsAnonymous
	note.Title = req.Title
	note.Content = req.Content

	err = dal.UpdateNote(note)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.ModifyNoteRes{}
}

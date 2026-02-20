package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

func DeleteNote(ctx *mhttp.Context) {
	req := &api.DeleteNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	note, err := dal.GetNote(req.NoteID)
	if err != nil {
		ctx.ResData = err
		return
	}

	if ctx.UserName != note.Writer {
		e := ErrNeedOwner().WithParam("operator", ctx.UserName).WithParam("owner", note.Writer)
		ctx.ResData = e
		mlog.Error(e.String())
		return
	}

	err = dal.DeleteNote(req.NoteID)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.DeleteNoteRes{}
}

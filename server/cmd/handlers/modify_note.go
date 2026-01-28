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

	note, err := dal.GetNote(req.ID)
	if err != nil {
		ctx.ResData = err
		return
	}

	if req.IsAnonymous == note.IsAnonymous && req.Title == note.Title && req.Content == note.Content {
		e := NewError(ET_ParamsError, ED_NoChanges).WithParam("operator", ctx.UserID)
		ctx.ResData = e
		mlog.Log(e.String())
		return
	}

	if ctx.UserID != note.WriterID {
		e := NewError(ET_OperatorError, ED_NeedOwner).WithParam("operator", ctx.UserID).WithParam("owner", note.WriterID)
		ctx.ResData = e
		mlog.Log(e.String())
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

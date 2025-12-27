package handlers

import (
	. "github.com/mats0319/unnamed_plan/server/internal/const"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
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

	if ctx.UserID != note.WriterID {
		e := NewError(ET_OperatorError, ED_NeedOwner).WithParam("operator", ctx.UserID).WithParam("owner", note.WriterID)
		mlog.Log(e.String())
		ctx.ResData = e
		return
	}
	if note.IsAnonymous == req.IsAnonymous && note.Title == req.Title && note.Content == req.Content {
		e := NewError(ET_ParamsError, ED_ModifyNothing).WithParam("operator", ctx.UserID)
		mlog.Log(e.String())
		ctx.ResData = e
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

	ctx.ResData = &api.ModifyNoteRes{ResBase: api.ResBaseSuccess}
}

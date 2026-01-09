package handlers

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

func DeleteNote(ctx *mhttp.Context) {
	req := &api.DeleteNoteReq{}
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
		ctx.ResData = e
		mlog.Log(e.String())
		return
	}

	err = dal.DeleteNote(req.ID)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.DeleteNoteRes{ResBase: api.ResBaseSuccess}
}

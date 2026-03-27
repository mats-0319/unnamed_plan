package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
)

func CreateNote(ctx *mhttp.Context) {
	req := &api.CreateNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	operator, e := dal.GetUser(ctx.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	note := model.NewNote(operator.UserName, operator.Nickname, req.IsAnonymous, req.Title, req.Content)
	if e := dal.CreateNote(note); e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.CreateNoteRes{}
}

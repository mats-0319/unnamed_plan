package handlers

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func CreateNote(ctx *mhttp.Context) {
	req := &api.CreateNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	operator, err := dal.GetUser(ctx.UserID)
	if err != nil {
		ctx.ResData = err
		return
	}

	note := model.NewNote(operator.ID, operator.Nickname, req.IsAnonymous, req.Title, req.Content)
	err = dal.CreateNote(note)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.CreateNoteRes{}
}

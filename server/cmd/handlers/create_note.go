package handlers

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
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

	note := &model.Note{
		WriterID:    operator.ID,
		WriterName:  operator.Nickname,
		IsAnonymous: req.IsAnonymous,
		Title:       req.Title,
		Content:     req.Content,
	}
	note.NoteID = utils.CalcSHA256(note.Serialize()) // 保证接口幂等性
	err = dal.CreateNote(note)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.CreateNoteRes{ResBase: api.ResBaseSuccess}
}

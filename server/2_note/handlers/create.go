package handlers

import (
	"github.com/mats0319/unnamed_plan/server/2_note/config"
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

	reader, err := ctx.Forward("http://" + config.ConfigIns.Address + "/api" + api.URI_Authenticate)
	if err != nil {
		ctx.ResData = err
		return
	}

	auth := &api.AuthenticateRes{}
	if !ctx.ParseParams(auth, reader) {
		return
	}

	note := &model.Note{
		NoteID:      utils.Uuid[string](),
		WriterID:    auth.UserID,
		WriterName:  auth.Nickname,
		IsAnonymous: req.IsAnonymous,
		Title:       req.Title,
		Content:     req.Content,
	}
	if err = dal.CreateNote(note); err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.CreateNoteRes{ResBase: api.ResBaseSuccess}
}

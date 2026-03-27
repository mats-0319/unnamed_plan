package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListNote(ctx *mhttp.Context) {
	req := &api.ListNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	var (
		count int64
		notes = make([]*model.Note, 0)
		e     *utils.Error
	)
	if req.OnlyOperator && len(ctx.UserName) > 0 {
		count, notes, e = dal.ListNote(req.Page.Size, req.Page.Num, ctx.UserName)
	} else {
		count, notes, e = dal.ListNote(req.Page.Size, req.Page.Num, "")
	}
	if e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ListNoteRes{
		Amount: count,
		Notes:  notesFromDBToHttp(notes),
	}
}

func notesFromDBToHttp(notes []*model.Note) []*api.Note {
	res := make([]*api.Note, len(notes))
	for i, v := range notes {
		res[i] = &api.Note{
			NoteID:      v.NoteID,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			Writer:      v.WriterName,
			IsAnonymous: v.IsAnonymous,
			Title:       v.Title,
			Content:     v.Content,
		}
	}

	return res
}

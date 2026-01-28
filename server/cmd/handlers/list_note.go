package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/model"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
)

func ListNote(ctx *mhttp.Context) {
	req := &api.ListNoteReq{}
	if !ctx.ParseParams(req) {
		return
	}

	count, notes, err := dal.ListNote(req.Page, req.UserID)
	if err != nil {
		ctx.ResData = err
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
			ID:          v.ID,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			NoteID:      v.NoteID,
			Writer:      v.WriterName,
			IsAnonymous: v.IsAnonymous,
			Title:       v.Title,
			Content:     v.Content,
		}
	}

	return res
}

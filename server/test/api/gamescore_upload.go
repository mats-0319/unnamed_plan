package api

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func UploadGameScore() {
	testCase("invalid game name", uploadGameScoreCase_InvalidGameName)
	testCase("success - visitor", uploadGameScoreCase_SuccessVisitor)
	testCase("success - user", uploadGameScoreCase_SuccessUser)
}

func uploadGameScoreCase_InvalidGameName() string {
	res := httpInvoke(api.URI_UploadGameScore, `{"game_name":-1,"score":100000,"result":"test game result","player":"Visotor001"}`, "")
	if res.IsSuccess || res.Err != utils.ErrInvalidGameName().Error() {
		return unknownError
	}

	return ""
}

func uploadGameScoreCase_SuccessVisitor() string {
	res := httpInvoke(api.URI_UploadGameScore, `{"game_name":1,"score":100000,"result":"test game result","player":"Visotor001"}`, "")
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListFlipGameScore(10, 1)
	if err != nil || count != 1 {
		return unknownError
	}

	return ""
}

func uploadGameScoreCase_SuccessUser() string {
	res := httpInvoke(api.URI_UploadGameScore, `{"game_name":1,"score":100000,"result":"test game result","player":""}`, accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListFlipGameScore(10, 1)
	if err != nil || count != 2 {
		return unknownError
	}

	return ""
}

package api

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListGameScore() {
	testCase("invalid game name", listGameScoreCase_InvalidGameName)
	testCase("success", listGameScoreCase_Success)
}

func listGameScoreCase_InvalidGameName() string {
	res := httpInvoke(api.URI_ListGameScore, `{"game_name":-1,"page":{"size":10,"num":1}}`, "")
	if res.IsSuccess || res.Err != utils.ErrInvalidGameName().Error() {
		return unknownError
	}

	return ""
}

func listGameScoreCase_Success() string {
	res := httpInvoke(api.URI_ListGameScore, `{"game_name":1,"page":{"size":10,"num":1}}`, "")
	if !res.IsSuccess {
		return res.Err
	}

	return ""
}

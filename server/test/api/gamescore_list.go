package api

import api "github.com/mats0319/unnamed_plan/server/cmd/api/go"

func ListGameScore() {
	testCase("success", listGameScoreCase_Success)
}

func listGameScoreCase_Success() string {
	res := httpInvoke(api.URI_ListGameScore, `{"game_name":1,"page":{"size":10,"num":1}}`, "")
	if !res.IsSuccess {
		return res.Err
	}

	return ""
}

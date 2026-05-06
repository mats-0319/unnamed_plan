package main

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

var defaultUser = []*model.User{
	newUser("mats0319", "Mario", true, false, ""),
}

func newUser(userName string, nickname string, isAdmin bool, enable2FA bool, totpKey string) *model.User {
	pwdSHA256 := utils.CalcSHA256("123456")
	pwdArgon2 := password.GeneratePassword(pwdSHA256)

	if len(nickname) < 1 {
		nickname = userName
	}

	return &model.User{
		UserName:  userName,
		Nickname:  nickname,
		Password:  pwdArgon2,
		IsAdmin:   isAdmin,
		EnableMFA: enable2FA,
		TOTPKey:   totpKey,
	}
}

var testFlipGameScore = []*model.FlipGameScore{{
	GameScore: model.GameScore{
		Score:      10000,
		Result:     "test result",
		Player:     "",
		PlayerName: "visitor 001",
	},
}}

var testUser = []*model.User{
	newUser("admin", "", true, true, "5SSFNNEJUENPCCKP"),
	newUser("user", "", false, false, ""),
}

var testNote = []*model.Note{
	{
		Writer:      "mats0319",
		WriterName:  "Mario",
		IsAnonymous: false,
		Title:       "test title 1",
		Content:     "test content 1",
	},
	{
		Writer:      "mats0319",
		WriterName:  "Mario",
		IsAnonymous: false,
		Title:       "test title 2",
		Content:     "test content 2",
	},
}

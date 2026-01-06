package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

type LoginToken struct {
	UserID string
	Token  string
}

var token = &LoginToken{}

func HttpInvoke(uri string, payload string, t ...*LoginToken) string {
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:10319/api"+uri, strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	// set login token
	if len(t) > 0 {
		req.Header.Set(mconst.HttpHeader_UserID, t[0].UserID)
		req.Header.Set(mconst.HttpHeader_AccessToken, t[0].Token)
	} else {
		req.Header.Set(mconst.HttpHeader_UserID, token.UserID)
		req.Header.Set(mconst.HttpHeader_AccessToken, token.Token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := &api.ResBase{} // 因为会设计失败的用例，所以不判断r.isSuccess
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		log.Fatal(err)
	}
	if !r.IsSuccess {
		log.Println(r.Err)
	}

	// read login token
	userIDStr := res.Header.Get(mconst.HttpHeader_UserID)
	tokenStr := res.Header.Get(mconst.HttpHeader_AccessToken)

	if len(userIDStr) > 0 && len(tokenStr) > 0 {
		token.UserID = userIDStr
		token.Token = tokenStr
	}

	return r.Err
}

func TestApi(name string) {
	log.Printf("> %s.\n", name)
}

func TestCase(name string) {
	log.Printf("- %s:\n", name)
}

func TestApiEnd() {
	log.Println()
}

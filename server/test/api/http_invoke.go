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

var accessToken = ""

func HttpInvoke(uri string, payload string) string {
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:10319/api"+uri, strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set(mconst.HttpHeader_AccessToken, accessToken)

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

	// read access token
	tokenStr := res.Header.Get(mconst.HttpHeader_AccessToken)

	if len(tokenStr) > 0 {
		accessToken = tokenStr
	}

	return string(bodyBytes)
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

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

var accessToken_Admin = ""
var accessToken_User = ""

var pwd = utils.CalcSHA256("123456")

func GetAccessToken() {
	r := httpInvoke(api.URI_Login, fmt.Sprintf(`{"user_name":"admin","password":"%s","totp_code":""}`, pwd), "")
	if r.IsSuccess {
		accessToken_Admin = r.AccessToken
	}

	r = httpInvoke(api.URI_Login, fmt.Sprintf(`{"user_name":"user","password":"%s","totp_code":""}`, pwd), "")
	if r.IsSuccess {
		accessToken_User = r.AccessToken
	}
}

type TestResponse struct {
	IsSuccess   bool   `json:"is_success"`
	Err         string `json:"err"`
	AccessToken string `json:"access_token"`
}

func httpInvoke(uri string, payload string, token string) *TestResponse {
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:10319/api"+uri, strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set(utils.HttpHeader_AccessToken, token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := &TestResponse{}
	if err := json.Unmarshal(bodyBytes, r); err != nil {
		log.Fatal(err)
	}
	if r.IsSuccess {
		// read access token
		tokenStr := res.Header.Get(utils.HttpHeader_AccessToken)

		if len(tokenStr) > 0 {
			r.AccessToken = tokenStr
		}
	}

	return r
}

const unknownError = "unknown error"

func testCase(name string, f func() string) {
	if errStr := f(); len(errStr) < 1 {
		log.Printf("- case: %s. Test Passed.\n", name)
	} else {
		log.Printf("- case: %s. Test Failed, error: %s\n", name, errStr)
		panic("test failed")
	}
}

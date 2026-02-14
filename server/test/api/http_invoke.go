package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mconst "github.com/mats0319/unnamed_plan/server/internal/utils"
)

var accessToken = ""

func httpInvoke(uri string, payload string) *mhttp.Response {
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

	r := &mhttp.Response{} // 因为会设计失败的用例，所以不判断r.isSuccess
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		log.Fatal(err)
	}
	if r.IsSuccess {
		// read access token
		tokenStr := res.Header.Get(mconst.HttpHeader_AccessToken)

		if len(tokenStr) > 0 {
			accessToken = tokenStr
		}
	}

	return r
}

const unknownError = "unknown error"

func testCase(name string, f func() string) {
	errStr := f()
	if len(errStr) < 1 {
		log.Printf("- case: %s. Test Passed.\n", name)
	} else {
		log.Printf("- case: %s. Test Failed. error: %s\n", name, errStr)
		panic("test failed")
	}
}

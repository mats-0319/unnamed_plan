package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test(t *testing.T) {

	url := "http://127.0.0.1:10319/api/user/list"
	method := "POST"

	payload := strings.NewReader(`{
    "page": {
        "size": 10,
        "num": 1
    }
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Unnamed-Plan-User-ID", "1")
	req.Header.Set("Unnamed-Plan-Access-Token", "GhLOJFslvh")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

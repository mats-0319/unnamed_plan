package mhttp

import (
	"net/http"

	"github.com/mats0319/unnamed_plan/server/internal/const"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

func (ctx *Context) Invoke(url string) (*http.Response, *Error) {
	req, err := http.NewRequest(http.MethodPost, url, ctx.Request.Body)
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return nil, e
	}

	// set req header(s)
	req.Header.Add("Origin", ctx.Origin)
	for _, header := range mconst.HttpHeaderList {
		value := ctx.Request.Header.Get(header)
		req.Header.Add(header, value)
	}

	// 这里的err只是http连接出错一类的错误，http状态码不是200体现在res，err是nil
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return nil, e
	}

	// prepare res header(s)
	headers := make(map[string]string)
	for header, value := range res.Header {
		for _, v := range mconst.HttpHeaderList {
			if header == v {
				headers[header] = value[0]
				break
			}
		}
	}

	ctx.ResHeaders = headers

	return res, nil
}

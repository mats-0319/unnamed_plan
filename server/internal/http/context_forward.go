package mhttp

import (
	"bytes"
	"io"
	"net/http"

	. "github.com/mats0319/unnamed_plan/server/internal/const"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

func (ctx *Context) Forward(url string) (io.Reader, error) {
	req, err := http.NewRequest(http.MethodPost, url, ctx.Request.Body)
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return nil, e
	}

	req.Header.Add("Origin", ctx.Request.Host)
	for _, header := range HttpHeaderList { // forward our own http header
		value := ctx.Request.Header.Get(header)
		req.Header.Add(header, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return nil, e
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return nil, e
	}
	reader := bytes.NewReader(resBytes)

	resBaseIns := &api.ResBase{}
	if !ctx.ParseParams(resBaseIns, reader) || !resBaseIns.IsSuccess {
		e := NewError(ET_ServerInternalError).WithParam("res", *resBaseIns)
		mlog.Log(e.String())
		return nil, e
	}
	_, _ = reader.Seek(0, 0)

	// set our own header
	headers := make(map[string]string)
	for header, value := range res.Header {
		for _, v := range HttpHeaderList {
			if header == v {
				headers[header] = value[0]
				break
			}
		}
	}

	ctx.ResHeaders = headers

	return reader, nil
}

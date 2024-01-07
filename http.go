package easierweb

import (
	"bytes"
	"io"
	"net/http"
)

const (
	MethodGET     = "GET"
	MethodHEAD    = "HEAD"
	MethodOPTIONS = "OPTIONS"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
	MethodPATCH   = "PATCH"
	MethodDELETE  = "DELETE"
)

func GET(url string, header ...Params) (int, Data, error) {
	return HTTP(MethodGET, url, nil, header...)
}

func HEAD(url string, header ...Params) (int, Data, error) {
	return HTTP(MethodHEAD, url, nil, header...)
}

func OPTIONS(url string, header ...Params) (int, Data, error) {
	return HTTP(MethodOPTIONS, url, nil, header...)
}

func POST(url string, body Data, header ...Params) (int, Data, error) {
	return HTTP(MethodPOST, url, body, header...)
}

func PUT(url string, body Data, header ...Params) (int, Data, error) {
	return HTTP(MethodPUT, url, body, header...)
}

func PATCH(url string, body Data, header ...Params) (int, Data, error) {
	return HTTP(MethodPATCH, url, body, header...)
}

func DELETE(url string, header ...Params) (int, Data, error) {
	return HTTP(MethodDELETE, url, nil, header...)
}

func HTTP(method, url string, body Data, header ...Params) (int, Data, error) {
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	if len(header) > 0 {
		for _, h := range header {
			for k, v := range h {
				request.Header.Set(k, v)
			}
		}
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}
	return response.StatusCode, result, nil
}

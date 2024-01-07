package utils

import (
	"bytes"
	"github.com/dpwgc/easierweb"
	"io"
	"net/http"
)

func GET(url string, header ...map[string]string) (int, easierweb.Data, error) {
	return HTTP(easierweb.MethodGET, url, nil, header...)
}

func HEAD(url string, header ...map[string]string) (int, easierweb.Data, error) {
	return HTTP(easierweb.MethodHEAD, url, nil, header...)
}

func OPTIONS(url string, header ...map[string]string) (int, easierweb.Data, error) {
	return HTTP(easierweb.MethodOPTIONS, url, nil, header...)
}

func POST(url string, body []byte, header ...map[string]string) (int, easierweb.Data, error) {
	return HTTP(easierweb.MethodPOST, url, body, header...)
}

func PUT(url string, body []byte, header ...map[string]string) (int, easierweb.Data, error) {
	return HTTP(easierweb.MethodPUT, url, body, header...)
}

func PATCH(url string, body []byte, header ...map[string]string) (int, easierweb.Data, error) {
	return HTTP(easierweb.MethodPATCH, url, body, header...)
}

func DELETE(url string, header ...map[string]string) (int, easierweb.Data, error) {
	return HTTP(easierweb.MethodDELETE, url, nil, header...)
}

// HTTP return code, result and error
func HTTP(method, url string, body []byte, header ...map[string]string) (int, easierweb.Data, error) {
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

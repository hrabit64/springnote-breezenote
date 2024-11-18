package testUtils

import (
	"bytes"
	"encoding/json"
	"github.com/hrabit64/springnote-breezenote/core"
	"io"
	"net/http"
	"net/http/httptest"
)

// RunTestApiReq 테스트용 API 요청을 실행합니다.
func RunTestApiReq(method string, url string, body io.Reader, header http.Header) *httptest.ResponseRecorder {
	r := core.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	req.Header = header
	r.ServeHTTP(w, req)
	return w
}

// ConvertToJsonReader 데이터를 json 형식으로 변환하여 io.Reader로 반환합니다.
func ConvertToJsonReader(data interface{}) io.Reader {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(jsonData)
}

// ConvertResToMap 응답을 map으로 변환합니다.
func ConvertResToMap(res *httptest.ResponseRecorder) map[string]interface{} {
	var resMap map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resMap)
	if err != nil {
		panic(err)
	}

	return resMap
}

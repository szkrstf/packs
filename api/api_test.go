package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/szkrstf/packs/mock"
)

func TestCalculateHandler(t *testing.T) {
	h := CalculateHandler{calculator: &mock.Calculator{
		CalculateFn: func(i int) map[int]int { return map[int]int{10: 1} },
	}}

	tt := []struct {
		body   string
		status int
		resp   string
	}{
		{body: ``, status: 400, resp: "Bad request\n"},
		{body: `{"items":1`, status: 400, resp: "Bad request\n"},
		{body: `{"items":1}`, status: 200, resp: `{"data":{"10":1}}` + "\n"},
	}

	for _, tc := range tt {
		rw := &httptest.ResponseRecorder{Body: &bytes.Buffer{}}
		r, _ := http.NewRequest("POST", "", strings.NewReader(tc.body))
		r.Header.Add("Content-Type", "application/json")

		h.ServeHTTP(rw, r)

		if got, want := rw.Code, tc.status; got != want {
			t.Errorf("%v: got: %v; want: %v", tc.body, got, want)
		}
		if got, want := rw.Body.String(), tc.resp; got != want {
			t.Errorf("%v: got: %v; want: %v", tc.body, got, want)
		}
	}
}

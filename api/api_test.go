package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/szkrstf/packs"
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

func TestSizeHandler(t *testing.T) {

	tt := []struct {
		method string
		body   string
		get    func() []int
		set    func([]int) error

		status int
		resp   string
	}{
		{method: "GET", get: func() []int { return []int{} }, status: 200, resp: `{"data":[]}` + "\n"},
		{method: "GET", get: func() []int { return []int{1, 2, 3} }, status: 200, resp: `{"data":[1,2,3]}` + "\n"},
		{method: "POST", body: `{"sizes":[1,2,3]}`, set: func([]int) error { return nil }, status: 200},
		{method: "POST", body: `-{"sizes":[1,2,3]}`, status: 400, resp: "Bad request\n"},
		{method: "POST", body: `{"sizes":[1,2,3]}`, set: func([]int) error { return packs.ErrInvalidSizes }, status: 400, resp: "invalid sizes\n"},
		{method: "POST", body: `{"sizes":[1,2,3]}`, set: func([]int) error { return errors.New("internal error") }, status: 500, resp: "Internal error\n"},
		{method: "PUT", status: 405, resp: "Method not allowed\n"},
	}

	for _, tc := range tt {
		rw := &httptest.ResponseRecorder{Body: &bytes.Buffer{}}
		r, _ := http.NewRequest(tc.method, "", strings.NewReader(tc.body))
		r.Header.Add("Content-Type", "application/json")

		h := SizeHandler{store: &mock.SizeStore{
			GetFn: tc.get,
			SetFn: tc.set,
		}}

		h.ServeHTTP(rw, r)

		if got, want := rw.Code, tc.status; got != want {
			t.Errorf("%v: got: %v; want: %v", tc.body, got, want)
		}
		if got, want := rw.Body.String(), tc.resp; got != want {
			t.Errorf("%v: got: %v; want: %v", tc.body, got, want)
		}
	}
}

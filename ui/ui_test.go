package ui

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/szkrstf/packs/mock"
)

func TestCalculateHandler(t *testing.T) {
	srv := &mock.Calculator{
		CalculateFn: func(i int) map[int]int { return map[int]int{10: 1} },
	}
	h := NewCalculateHandler(srv)

	tt := []struct {
		form   map[string]string
		status int
		resp   []string
	}{
		{form: nil, status: 400, resp: []string{"Bad request"}},
		{form: map[string]string{"items": "-1"}, status: 400, resp: []string{"Bad request"}},
		{form: map[string]string{"invalid": "1"}, status: 400, resp: []string{"Bad request"}},
		{form: map[string]string{"items": "1"}, status: 200, resp: []string{"value=\"1\"", "<span>10: 1</span>"}},
	}

	for _, tc := range tt {
		rw := &httptest.ResponseRecorder{Body: &bytes.Buffer{}}
		form := make(url.Values)
		for k, v := range tc.form {
			form.Add(k, v)
		}
		r, _ := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		h.ServeHTTP(rw, r)

		if got, want := rw.Code, tc.status; got != want {
			t.Errorf("%v: got: %v; want: %v", tc.form, got, want)
		}
		for _, s := range tc.resp {
			if !strings.Contains(rw.Body.String(), s) {
				t.Errorf("%v: body doesn't contain %s", tc.form, s)
			}
		}
	}
}

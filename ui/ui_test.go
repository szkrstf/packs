package ui

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/szkrstf/packs"
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

func TestSizeHandler(t *testing.T) {
	tt := []struct {
		form   map[string]string
		get    func() []int
		set    func([]int) error
		status int
	}{
		{form: map[string]string{"invalid": "-1\r\n2\r\n3\r\n"}, set: func([]int) error { return packs.ErrInvalidSizes }, status: 400},
		{form: map[string]string{"sizes": "invalid\r\n2\r\n3\r\n"}, status: 400},
		{form: map[string]string{"sizes": "1\r\n2\r\n3\r\n"}, set: func([]int) error { return fmt.Errorf("internal error") }, status: 500},
		{form: map[string]string{"sizes": "1\r\n2\r\n3\r\n"}, set: func([]int) error { return nil }, status: 200},
	}

	for _, tc := range tt {
		rw := &httptest.ResponseRecorder{Body: &bytes.Buffer{}}
		form := make(url.Values)
		for k, v := range tc.form {
			form.Add(k, v)
		}
		r, _ := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		h := NewSizeHandler(&mock.SizeStore{
			GetFn: func() []int { return []int{1, 2, 3} },
			SetFn: tc.set,
		})

		h.ServeHTTP(rw, r)

		if got, want := rw.Code, tc.status; got != want {
			t.Errorf("%v: got: %v; want: %v", tc.form, got, want)
		}
	}
}

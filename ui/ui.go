package ui

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/szkrstf/packs"
)

// CalculateHandler is a http.Handler for the ui.
type CalculateHandler struct {
	calculator packs.Calculator

	calcTmpl *template.Template
}

// NewCalculateHandler creates a new ui.CalculateHandler.
func NewCalculateHandler(calculator packs.Calculator) *CalculateHandler {
	return &CalculateHandler{
		calculator: calculator,
		calcTmpl:   template.Must(template.New("calc").Parse(calcHTML)),
	}
}

type calcResp struct {
	Items int
	Data  map[int]int
}

// ServeHTTP implements the http.Handler interface.
func (h *CalculateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := h.calcTmpl.Execute(w, nil); err != nil {
			log.Println(err)
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	items, err := strconv.Atoi(r.FormValue("items"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if items <= 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	packs := h.calculator.Calculate(items)

	if err := h.calcTmpl.Execute(w, calcResp{Items: items, Data: packs}); err != nil {
		log.Println(err)
	}
}

const calcHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>packs</title>
	</head>
	<body>
		<p>
			<a href="/sizes">edit sizes</a>
		</p>
		<form method="POST">
			<input type="text" name="items" value="{{ .Items }}" />
			<input type="submit" value="Calculate" />
		</form>
		<ul>
			{{- range $k, $v := .Data }}
			<li><span>{{ $k }}: {{ $v }}</span></li>
			{{- end }}
		</ul>
	</body>
</html>
`

// SizeHandler is a http.Handler for the ui.
type SizeHandler struct {
	store packs.SizeStore

	sizeTmpl *template.Template
}

// NewSizeHandler creates a new ui.Handler.
func NewSizeHandler(store packs.SizeStore) *SizeHandler {
	return &SizeHandler{
		store:    store,
		sizeTmpl: template.Must(template.New("size").Parse(sizeHTML)),
	}
}

type sizeResp struct {
	Data string
}

// ServeHTTP implements the http.Handler interface.
func (h *SizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		sizes := h.store.Get()
		if err := h.sizeTmpl.Execute(w, sizeResp{Data: printSizes(sizes)}); err != nil {
			log.Println(err)
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	sizes, err := readSizes(strings.NewReader(r.FormValue("sizes")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.Set(sizes); errors.Is(err, packs.ErrInvalidSizes) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if err := h.sizeTmpl.Execute(w, sizeResp{Data: printSizes(sizes)}); err != nil {
		log.Println(err)
	}
}

func readSizes(r io.Reader) ([]int, error) {
	var sizes []int
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		s, err := strconv.Atoi(strings.TrimSpace(sc.Text()))
		if err != nil {
			return nil, fmt.Errorf("invalid sizes: not a number")
		}
		sizes = append(sizes, s)
	}
	return sizes, nil
}

func printSizes(sizes []int) string {
	var res string
	for _, s := range sizes {
		res += fmt.Sprintf("%d\n", s)
	}
	return res
}

const sizeHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>packs</title>
	</head>
	<body>
		<p>
			<a href="/">calculate packs</a>
		</p>
		<form method="POST" action="/sizes">
			<textarea name="sizes">{{ .Data }}</textarea>
			<input type="submit" value="Save" />
		</form>
	</body>
</html>
`

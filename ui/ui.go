package ui

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

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

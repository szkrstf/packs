package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/szkrstf/packs"
	"github.com/szkrstf/packs/api"
)

func main() {
	addr := flag.String("addr", ":8080", "http address")
	flag.Parse()

	sizes := []int{250, 500, 1000, 2000, 5000}
	calculator, err := packs.NewCalculator(sizes)
	if err != nil {
		log.Panic(err)
	}

	http.Handle("/api/calculate", api.NewCalculateHangler(calculator))

	log.Println("server listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Panic(err)
	}
}

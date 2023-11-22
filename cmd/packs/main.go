package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/szkrstf/packs"
	"github.com/szkrstf/packs/api"
	"github.com/szkrstf/packs/ui"
)

func main() {
	addr := flag.String("addr", ":8080", "http address")
	configFile := flag.String("config", "", "config file location for sizes")
	flag.Parse()

	sizes := []int{250, 500, 1000, 2000, 5000}
	if *configFile != "" {
		var err error
		sizes, err = readConfigFile(*configFile)
		if err != nil {
			log.Fatalf("error reading config(%s): %v", *configFile, err)
		}
	}

	calculator, err := packs.NewCalculator(sizes)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/calculate", api.NewCalculateHangler(calculator))
	http.Handle("/", ui.NewCalculateHandler(calculator))

	log.Println("server listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func readConfigFile(name string) ([]int, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sizes []int
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s, err := strconv.Atoi(strings.TrimSpace(sc.Text()))
		if err != nil {
			return nil, fmt.Errorf("invalid config file: %v", err)
		}
		sizes = append(sizes, s)
	}
	return sizes, nil
}

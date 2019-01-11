package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/dhlk/exrays"
)

var (
	address string
)

func flagInit() {
	flag.StringVar(&address, "a", ":5475", "listen address")

	flag.Parse()
}

type AppImgSort []exrays.AppImg

func (ais AppImgSort) Len() int {
	return len(ais)
}

func (ais AppImgSort) Less(i, j int) bool {
	return ais[i].Time.After(ais[j].Time)
}

func (ais AppImgSort) Swap(i, j int) {
	ais[i], ais[j] = ais[j], ais[i]
}

func pullApps(cutoff time.Time, apps []string) []exrays.AppImg {
	results := make([]exrays.AppImg, 0)

	for _, app := range apps {
		for _, result := range exrays.PullApp(app) {
			if result.Time.After(cutoff) {
				results = append(results, result)
			}
		}
	}

	sort.Sort(AppImgSort(results))

	return results
}

func writeImgs(imgs []exrays.AppImg, w io.Writer, t exrays.Transform) {
	for _, img := range imgs {
		fmt.Fprintf(w, `<a href="`+img.Link+`">`)
		fmt.Fprintf(w, `<img src="data:image/png;base64,`)

		resp, err := http.Get(img.Image)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		b64w := base64.NewEncoder(base64.StdEncoding, w)
		exrays.Decode(resp.Body, b64w, t)
		b64w.Close()

		fmt.Fprintf(w, `" /></a><br>`)
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	req.Body.Close()
	req.ParseForm()

	d := req.Form.Get("d")
	if d == "" {
		d = "1m"
	}
	dur, err := time.ParseDuration(d)
	if err != nil {
		panic(err)
	}
	cutoff := time.Now().Add(-1 * dur)

	t := req.Form.Get("t")
	if t == "" {
		t = "r"
	}

	fmt.Fprintf(w, "<html><body>")
	writeImgs(pullApps(cutoff, req.Form["a"]), w, exrays.Transforms[t])
	fmt.Fprintf(w, "<h1>done</h1></body></html>")
}

func main() {
	flagInit()
	http.HandleFunc("/", handler)

	for {
		log.Printf("%v", http.ListenAndServe(address, nil))
	}
}

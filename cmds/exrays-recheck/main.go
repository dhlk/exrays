package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	//"path"
	//"strings"
	"sync"
	"time"

	"github.com/dhlk/exrays"
)

var (
	address string
	apps    []string
)

func flagInit() []string {
	flag.StringVar(&address, "a", ":5475", "listen address")

	flag.Parse()
	return flag.Args()
}

func feeder(app string, imgs chan<- exrays.AppImg, wait sync.WaitGroup, cutoff time.Time) {
	defer wait.Done()
	results := exrays.PullApp(app)

	for _, result := range results {
		if result.Time.Before(cutoff) {
			continue
		}

		//if strings.ToLower(path.Ext(result.Image)) != ".png" {
		//	continue
		//}

		imgs <- result
	}
}

func writer(imgs <-chan exrays.AppImg, w io.Writer, t exrays.Transform) {
	for img := range imgs {
		if img.Image == "" {
			break
		}

		fmt.Fprintf(w, `<img src="data:image/png;base64,`)

		resp, err := http.Get(img.Image)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		b64w := base64.NewEncoder(base64.StdEncoding, w)
		exrays.Decode(resp.Body, b64w, t)
		b64w.Close()

		fmt.Fprintf(w, `" />`)
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

	imgs := make(chan exrays.AppImg)
	go writer(imgs, w, exrays.Transforms[t])

	var wait sync.WaitGroup
	wait.Add(len(apps))
	for _, app := range apps {
		go feeder(app, imgs, wait, cutoff)
	}

	wait.Wait()
	close(imgs)

	fmt.Fprintf(w, "</body></html>")
}

func main() {
	apps = flagInit()
	http.HandleFunc("/", handler)

	for {
		log.Printf("%v", http.ListenAndServe(address, nil))
	}
}

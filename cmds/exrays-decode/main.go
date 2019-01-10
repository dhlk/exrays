package main

import (
	"flag"
	"log"
	"os"

	"github.com/dhlk/exrays"
)

var (
	tx string
)

func initFlags() []string {
	flag.StringVar(&tx, "t", "r", "transform")
	flag.Parse()
	return flag.Args()
}

func decode(source, destination string) {
	r, err := os.Open(source)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
	defer r.Close()

	w, err := os.Create(destination)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
	defer w.Close()

	exrays.Decode(r, w, exrays.Transforms[tx])
	log.Printf("%s => %s", source, destination)
}

func main() {
	args := initFlags()
	nargs := len(args)

	for i := 0; i < (nargs/2)*2; i += 2 {
		decode(args[i], args[i+1])
	}
}

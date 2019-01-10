package main

import (
	"log"
	"os"

	"github.com/dhlk/exrays"
)

func printImgs(imgs []exrays.AppImg) {
	for _, img := range imgs {
		log.Printf("%v", img)
	}
}

func main() {
	for _, arg := range os.Args[1:] {
		printImgs(exrays.PullApp(arg))
	}
}

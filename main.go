package main

import (
	"fmt"

	"github.com/seballgeyer/goginan/gpxparser"
)

func main() {
	fileName := "scratch/DARW-2019-07-18_00:00.gpx"
	data := gpxparser.ParseGpx(fileName)
	fmt.Println(data)
}

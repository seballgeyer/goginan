package main

import (
	"github.com/seballgeyer/goginan/gpxparser"
	"github.com/seballgeyer/goginan/position"
)

func main() {
	fileName := "scratch/DARW-2019-07-18_00:00.gpx"
	data := gpxparser.ParseGpx(fileName)
	position.FillEmptyPosData(&data)
	data.Output("test.txt")
	// for _, d := range data.Data {
	// 	fmt.Println(d)
	// }
	// fmt.Printf("Station: %s, apriori %v \n", data.Station, data.Apriori)
}

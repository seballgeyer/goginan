package gpxparser

import (
	"encoding/xml"
	"io"
	"log"
	"os"

	"github.com/seballgeyer/goginan/position"
)

type gpx struct {
	// XMLName xml.Name `xml:"gpx"`
	Trk trk `xml:"trk"`
}

type trk struct {
	// XMLName xml.Name `xml:"trk"`
	Trkseg []trkseg `xml:"trkseg"`
}

type trkseg struct {
	// XMLName xml.Name `xml:"trkseg"`
	Trkpt []trkpt `xml:"trkpt"`
}

type trkpt struct {
	// XMLName    xml.Name   `xml:"trkpt"`
	Lat        float64    `xml:"lat,attr"`
	Lon        float64    `xml:"lon,attr"`
	Ele        float64    `xml:"ele"`
	Time       string     `xml:"time"`
	Extensions extensions `xml:"extensions"`
}

type extensions struct {
	// XMLName xml.Name `xml:"extensions"`
	Time    string `xml:"time"`
	Pos     XYZ    `xml:"pos"`
	Vcv     VCV    `xml:"variances"`
	Apriori XYZ    `xml:"apriori"`
}

type XYZ struct {
	// XMLName xml.Name `xml:"pos"`
	X float64 `xml:"x"`
	Y float64 `xml:"y"`
	Z float64 `xml:"z"`
}

type VCV struct {
	// XMLName xml.Name `xml:"variances"`
	Xx float64 `xml:"xx"`
	Xy float64 `xml:"xy"`
	Xz float64 `xml:"xz"`
	Yy float64 `xml:"yy"`
	Yz float64 `xml:"yz"`
	Zz float64 `xml:"zz"`
}

func ParseGpx(fileName string) position.Position {
	xmlFile, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	var gpx gpx
	xml.Unmarshal(byteValue, &gpx)

	apr := gpx.Trk.Trkseg[0].Trkpt[0].Extensions.Apriori
	var data position.Position
	for _, trkseg := range gpx.Trk.Trkseg {
		for _, trkpt := range trkseg.Trkpt {
			pos := position.PositionData{
				Time: trkpt.Time,
				Lat:  trkpt.Lat,
				Lon:  trkpt.Lon,
				Elev: trkpt.Ele,
				Pos:  position.XYZ{X: trkpt.Extensions.Pos.X, Y: trkpt.Extensions.Pos.Y, Z: trkpt.Extensions.Pos.Z},
				Vcv:  position.VCV{Xx: trkpt.Extensions.Vcv.Xx, Xy: trkpt.Extensions.Vcv.Xy, Xz: trkpt.Extensions.Vcv.Xz, Yy: trkpt.Extensions.Vcv.Yy, Yz: trkpt.Extensions.Vcv.Yz, Zz: trkpt.Extensions.Vcv.Zz},
			}
			data.Data = append(data.Data, pos)

		}
	}
	data.Station = "DARW"
	data.Apriori = position.XYZ{X: apr.X, Y: apr.Y, Z: apr.Z}
	return data
}

package position

type XYZ struct {
	X float64
	Y float64
	Z float64
}

type VCV struct {
	Xx float64
	Xy float64
	Xz float64
	Yy float64
	Yz float64
	Zz float64
}

type Position struct {
	Time     string
	Lat      float64
	Lon      float64
	Elev     float64
	Pos      XYZ
	Vcv      VCV
	Apriori  XYZ
	DeltaPos XYZ
	DeltaENU XYZ
	VCVENU   VCV
}

// func fillEmptyPosData(data []Position) []Position {
// 	for i, d := range data {
// 		data[i].DeltaPos = d.Pos
// 	}
// 	return data
// }

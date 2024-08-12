package position

import (
	"fmt"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"
)

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

type PositionData struct {
	Time     string
	Lat      float64
	Lon      float64
	Elev     float64
	Pos      XYZ
	Vcv      VCV
	DeltaPos XYZ
	DeltaENU XYZ
	VCVENU   VCV
}

type Position struct {
	Station string
	Apriori XYZ
	Data    []PositionData
}

func FillEmptyPosData(data *Position) {
	_, rot := XYZToNEU(data.Apriori)
	fmt.Println(rot)
	for i, d := range data.Data {
		data.Data[i].DeltaPos = data.Apriori.Sub(d.Pos)
		data.Data[i].DeltaENU = data.Data[i].DeltaPos.Rotate(rot)
		data.Data[i].VCVENU = d.Vcv.Rotate(rot)
	}
	// return data
}

func (p XYZ) Sub(q XYZ) XYZ {
	return XYZ{X: p.X - q.X, Y: p.Y - q.Y, Z: p.Z - q.Z}
}

func XYZToNEU(site XYZ) (XYZ, mat.Dense) {
	// Convert ECEF to ENU

	// Rotation matrix
	// | -sin(lon)      cos(lon)       0 |
	// | -sin(lat)cos(lon) -sin(lat)sin(lon) cos(lat) |
	// |  cos(lat)cos(lon)  cos(lat)sin(lon) sin(lat) |
	lat := math.Atan2(site.Z, math.Sqrt(site.X*site.X+site.Y*site.Y))
	lon := math.Atan2(site.Y, site.X)
	rot := mat.NewDense(3, 3, []float64{
		-math.Sin(lon), math.Cos(lon), 0,
		-math.Sin(lat) * math.Cos(lon), -math.Sin(lat) * math.Sin(lon), math.Cos(lat),
		math.Cos(lat) * math.Cos(lon), math.Cos(lat) * math.Sin(lon), math.Sin(lat),
	})
	enu := XYZ{X: -site.Y, Y: -site.X, Z: site.Z}
	return enu, *rot
}

func (p *XYZ) Rotate(rot mat.Dense) XYZ {
	// Rotate the vector p by the rotation matrix rot
	// p = rot * p
	pVec := mat.NewVecDense(3, []float64{p.X, p.Y, p.Z})
	pVecRot := mat.NewVecDense(3, nil)
	pVecRot.MulVec(&rot, pVec)
	output := XYZ{X: pVecRot.AtVec(0), Y: pVecRot.AtVec(1), Z: pVecRot.AtVec(2)}
	return output
}

func (p *VCV) Rotate(rot mat.Dense) VCV {
	// Rotate the VCV p by the rotation matrix rot
	// p = rot * p * rot^T
	pMat := mat.NewDense(3, 3, []float64{
		p.Xx, p.Xy, p.Xz,
		p.Xy, p.Yy, p.Yz,
		p.Xz, p.Yz, p.Zz,
	})
	pMatRot := mat.NewDense(3, 3, nil)
	pMatRot.Mul(&rot, pMat)
	pMatRot.Mul(pMatRot, rot.T())
	output := VCV{
		Xx: pMatRot.At(0, 0),
		Xy: pMatRot.At(0, 1),
		Xz: pMatRot.At(0, 2),
		Yy: pMatRot.At(1, 1),
		Yz: pMatRot.At(1, 2),
		Zz: pMatRot.At(2, 2),
	}
	return output
}

func (p Position) Output(filename string) error {
	//open file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	//write to file
	for _, d := range p.Data {
		_, err := file.WriteString(fmt.Sprintf("%s % 14.10f % 14.10f % 10.5f ", d.Time, d.Lat, d.Lon, d.Elev))
		if err != nil {
			return err
		}

		_, err = file.WriteString(fmt.Sprintf("% 13.5f % 13.5f % 13.5f ", d.Pos.X, d.Pos.Y, d.Pos.Z))
		if err != nil {
			return err
		}

		_, err = file.WriteString(fmt.Sprintf("% 7.5f % 7.5f % 7.5f % 7.5f % 7.5f % 7.5f ", d.Vcv.Xx, d.Vcv.Yy, d.Vcv.Zz, d.Vcv.Xy, d.Vcv.Xz, d.Vcv.Yz))
		if err != nil {
			return err
		}

		_, err = file.WriteString(fmt.Sprintf("% 7.5f % 7.5f % 7.5f ", d.DeltaPos.X, d.DeltaPos.Y, d.DeltaPos.Z))
		if err != nil {
			return err
		}

		_, err = file.WriteString(fmt.Sprintf("% 7.5f % 7.5f % 7.5f ", d.DeltaENU.X, d.DeltaENU.Y, d.DeltaENU.Z))
		if err != nil {
			return err
		}

		_, err = file.WriteString(fmt.Sprintf("%10.7f %10.7f %10.7f %10.7f %10.7f %10.7f\n", d.VCVENU.Xx, d.VCVENU.Yy, d.VCVENU.Zz, d.VCVENU.Xy, d.VCVENU.Xz, d.VCVENU.Yz))
		if err != nil {
			return err
		}
	}

	// Output the Position data to a file
	return nil
	// return fmt.Sprintf("Time: %s, Lat: %f, Lon: %f, Elev: %f, Pos: %v, VCV: %v, DeltaPos: %v, DeltaENU: %v, VCVENU: %v", p.Time, p.Lat, p.Lon, p.Elev, p.Pos, p.Vcv, p.DeltaPos, p.DeltaENU, p.VCVENU)
}

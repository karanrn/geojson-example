package helper

import (
	"github.com/paulmach/orb"

	"github.com/karanrn/geojson-example/geodata"
)

// GetCentriod of the polygon
func GetCentriod(or orb.Ring) geodata.Loc {
	var centriod geodata.Loc
	n := len(or)
	signedArea := 0.0
	for i := 0; i < n; i++ {
		x0 := or[i][0]
		y0 := or[i][1]
		x1 := or[(i+1)%n][0]
		y1 := or[(i+1)%n][1]

		a := (x0 * y1) - (x1 * y0)
		signedArea += a

		centriod.Lat += (x0 + x1) * a
		centriod.Lon += (y0 + y1) * a
	}
	signedArea *= 0.5
	centriod.Lat = (centriod.Lat / (6 * signedArea))
	centriod.Lon = (centriod.Lon / (6 * signedArea))
	return centriod
}

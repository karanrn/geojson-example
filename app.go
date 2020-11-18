package main

/*
https://golangcode.com/is-point-within-polygon-from-geojson/
https://www.geeksforgeeks.org/find-the-centroid-of-a-non-self-intersecting-closed-polygon/

*/
import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

const (
	// GEOFILE is the GeoJSON file
	GEOFILE = "IndianStates.json"
)

/*
East to West (Lat) : 98 to 70
North to South (Lon): 36 to 8
*/
// Loc holds Latitude and Longitude information
type Loc struct {
	Lat float64
	Lon float64
}

// StateCentriod has centriod of the state
type StateCentriod struct {
	State    string
	Centriod Loc
}

type StateCentriodList []StateCentriod

func main() {
	// Open/load the file
	f, err := ioutil.ReadFile(GEOFILE)
	if err != nil {
		fmt.Errorf("error while reading json file, got %v", err.Error())
		return
	}

	featureCollections, err := geojson.UnmarshalFeatureCollection(f)
	/*
		// List all states
		for _, feature := range featureCollections.Features {
			fmt.Println(feature.Properties["NAME_1"])
		}
	*/

	/*
		// Find the state in which geolocation lies - [93.789047, 6.852571]
		result := isPointInsidePolygon(featureCollections, orb.Point{93.789047, 6.852571})
		if result == "" {
			fmt.Println("Given geolocation does not lie in the India.")
		} else {
			fmt.Println(result)
		}

	*/

	// Find centriods of the state
	var sCentriods []StateCentriod
	for _, feature := range featureCollections.Features {
		_, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			//fmt.Println(feature.Properties["NAME_1"], feature.Geometry.(orb.MultiPolygon)[0][0])
			sCentriods = append(sCentriods, StateCentriod{
				State:    feature.Properties["NAME_1"].(string),
				Centriod: getCentroid(feature.Geometry.(orb.MultiPolygon)[0][0]),
			})

		} else {
			sCentriods = append(sCentriods, StateCentriod{
				State:    feature.Properties["NAME_1"].(string),
				Centriod: getCentroid(feature.Geometry.(orb.Polygon)[0]),
			})
		}
		//break
	}

	// East to West, ordered alphabetically
	sort.Slice(sCentriods, func(i, j int) bool {
		if sCentriods[i].Centriod.Lat > sCentriods[j].Centriod.Lat {
			return true
		}
		if sCentriods[i].Centriod.Lat < sCentriods[j].Centriod.Lat {
			return false
		}
		return sCentriods[i].State < sCentriods[j].State
	})
	for _, s := range sCentriods {
		fmt.Println(s)
	}

	// North to South, ordered alphabetically
	sort.Slice(sCentriods, func(i, j int) bool {
		if sCentriods[i].Centriod.Lon > sCentriods[j].Centriod.Lon {
			return true
		}
		if sCentriods[i].Centriod.Lon < sCentriods[j].Centriod.Lon {
			return false
		}
		return sCentriods[i].State < sCentriods[j].State
	})
	for _, s := range sCentriods {
		fmt.Println(s)
	}
}

func getCentroid(or orb.Ring) Loc {
	var centriod Loc
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

func isPointInsidePolygon(fc *geojson.FeatureCollection, point orb.Point) string {
	for _, feature := range fc.Features {
		// Try on a MultiPolygon to begin
		multiPoly, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			if planar.MultiPolygonContains(multiPoly, point) {
				return feature.Properties["NAME_1"].(string)
			}
		} else {
			// Fallback to Polygon
			polygon, isPoly := feature.Geometry.(orb.Polygon)
			if isPoly {
				if planar.PolygonContains(polygon, point) {
					return feature.Properties["NAME_1"].(string)
				}
			}
		}
	}
	return ""
}

func (s StateCentriodList) Len() int {
	return len(s)
}

func (s StateCentriodList) Less(i, j int) bool {
	return s[i].Centriod.Lat > s[j].Centriod.Lat
}

func (s StateCentriodList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

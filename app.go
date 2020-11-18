package main

/*
https://golangcode.com/is-point-within-polygon-from-geojson/
https://www.geeksforgeeks.org/find-the-centroid-of-a-non-self-intersecting-closed-polygon/

*/
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

const (
	// GEOFILE is the GeoJSON file
	GEOFILE = "IndianStates.json"
)

// UTs is union territories of the India
var UTs = []string{"Andaman and Nicobar", "Chandigarh", "Dadra and Nagar Haveli", "Daman and Diu", "Delhi", "Jammu and Kashmir", "Ladakh", "Lakshadweep", "Puducherry"}

// Loc holds Latitude and Longitude information
/*
East to West (Lat) : 98 to 70
North to South (Lon): 36 to 8
*/
type Loc struct {
	Lat float64 `json:"Latitude"`
	Lon float64 `json:"Longitude"`
}

// StateCentriod has centriod of the state
type StateCentriod struct {
	State    string
	isUT     bool
	Centriod Loc
}

// StateCentriodList slice of type StateCentriod
type StateCentriodList []StateCentriod

var featureCollections geojson.FeatureCollection
var sCentriods []StateCentriod

func init() {
	DataPrep()
}

// DataPrep loads and prepares data for the API
func DataPrep() {
	// Open/load the file
	f, err := ioutil.ReadFile(GEOFILE)
	if err != nil {
		fmt.Errorf("error while reading json file, got %v", err.Error())
		return
	}

	featureCollections, err := geojson.UnmarshalFeatureCollection(f)

	// Find centriods of the state
	for _, feature := range featureCollections.Features {
		// Check if it is a UT
		ut := Contains(feature.Properties["NAME_1"].(string), UTs)

		_, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			sCentriods = append(sCentriods, StateCentriod{
				State:    feature.Properties["NAME_1"].(string),
				Centriod: getCentroid(feature.Geometry.(orb.MultiPolygon)[0][0]),
				isUT:     ut,
			})

		} else {
			sCentriods = append(sCentriods, StateCentriod{
				State:    feature.Properties["NAME_1"].(string),
				Centriod: getCentroid(feature.Geometry.(orb.Polygon)[0]),
				isUT:     ut,
			})
		}
	}
}

// ListStatesAndUT list states and union territories basis
func ListStatesAndUT(ut bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			var states []string
			for _, sc := range sCentriods {
				// Only states
				if !ut && !sc.isUT {
					states = append(states, sc.State)
				}
				// States and Union territories
				if ut {
					states = append(states, sc.State)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(states)
		} else {
			fmt.Fprintf(w, "Method not supported")
		}
	}
}

/*
// GetState returns state if geolocation point lies in it
func GetState(w http.ResponseWriter, r *http.Request) {
	var location Loc
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&location)
		if err != nil {
			json.NewEncoder(w).Encode(`{'error': 'Error in decoding JSON'}`)
			return
		}

		result := isPointInsidePolygon(featureCollections, orb.Point{93.789047, 6.852571})
		if result == "" {
			fmt.Println("Given geolocation does not lie in the India.")
		} else {
			fmt.Println(result)
		}

	}
}
*/

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/statesonly", ListStatesAndUT(false))
	mux.HandleFunc("/states-ut", ListStatesAndUT(true))
	fmt.Println("Serving on :9000")
	log.Fatal(http.ListenAndServe(":9000", mux))

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

	/*
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
	*/
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

// Contains checks if value exists in the list
func Contains(key string, utList []string) bool {
	for _, ut := range utList {
		if key == ut {
			return true
		}
	}
	return false
}

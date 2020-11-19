package geodata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

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
	IsUT     bool
	Centriod Loc
}

// StateCentriodList slice of type StateCentriod
type StateCentriodList []StateCentriod

// FeatureCollections hold geojson feature collection data
var FeatureCollections *geojson.FeatureCollection

// SCentriods holds centriods of the states
var SCentriods []StateCentriod

// ListStatesAndUT list states and union territories basis
func ListStatesAndUT(ut bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			var states []string
			for _, sc := range SCentriods {
				// Only states
				if !ut && !sc.IsUT {
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

// GetState returns state if geolocation point lies in it
func GetState(w http.ResponseWriter, r *http.Request) {
	var location Loc
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&location)
		if err != nil {
			json.NewEncoder(w).Encode(`{'error': 'Error in decoding JSON'}`)
			return
		}

		result := isPointInsidePolygon(FeatureCollections, orb.Point{location.Lat, location.Lon})
		if result == "" {
			fmt.Fprintf(w, "given geolocation does not lie in the India")
		} else {
			fmt.Fprintf(w, "(%f, %f) lies in (%s)", location.Lat, location.Lon, result)
		}
	} else {
		fmt.Fprintf(w, "Method not supported")
	}
}

// OrderStates returns states ordered basis direction
// WE for West to East, NS for North to South
func OrderStates(direction string, ut bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var states []string
		if r.Method == "GET" {
			// West to East
			if direction == "WE" {
				sort.Slice(SCentriods, func(i, j int) bool {
					return SCentriods[i].Centriod.Lat < SCentriods[j].Centriod.Lat
				})
			}

			// North to South
			if direction == "NS" {
				sort.Slice(SCentriods, func(i, j int) bool {
					return SCentriods[i].Centriod.Lon > SCentriods[j].Centriod.Lon
				})
			}

			for _, st := range SCentriods {
				// Only States
				if !ut && !st.IsUT {
					states = append(states, st.State)
				}
				// States and Union Territories
				if ut {
					states = append(states, st.State)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(states)
		} else {
			fmt.Fprintf(w, "Method not supported")
		}
	}
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

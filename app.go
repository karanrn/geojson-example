package main

/*
https://golangcode.com/is-point-within-polygon-from-geojson/
https://www.geeksforgeeks.org/find-the-centroid-of-a-non-self-intersecting-closed-polygon/

*/
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	"github.com/karanrn/geojson-example/geodata"
	"github.com/karanrn/geojson-example/helper"
)

const (
	// GEOFILE is the GeoJSON file
	GEOFILE = "IndianStates.json"
)

// UTs is union territories of the India
var UTs = []string{"Andaman and Nicobar", "Chandigarh", "Dadra and Nagar Haveli", "Daman and Diu", "Delhi", "Jammu and Kashmir", "Ladakh", "Lakshadweep", "Puducherry"}

func init() {
	DataPrep()
}

// DataPrep loads and prepares data for the API
func DataPrep() {
	// Open/load the file
	f, err := ioutil.ReadFile(GEOFILE)
	if err != nil {
		fmt.Printf("error while reading json file, got %v", err.Error())
		return
	}

	geodata.FeatureCollections, err = geojson.UnmarshalFeatureCollection(f)

	// Find centriods of the state
	for _, feature := range geodata.FeatureCollections.Features {
		// Check if it is a UT
		ut := helper.Contains(feature.Properties["NAME_1"].(string), UTs)

		_, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			geodata.SCentriods = append(geodata.SCentriods, geodata.StateCentriod{
				State:    feature.Properties["NAME_1"].(string),
				Centriod: helper.GetCentriod(feature.Geometry.(orb.MultiPolygon)[0][0]),
				IsUT:     ut,
			})

		} else {
			geodata.SCentriods = append(geodata.SCentriods, geodata.StateCentriod{
				State:    feature.Properties["NAME_1"].(string),
				Centriod: helper.GetCentriod(feature.Geometry.(orb.Polygon)[0]),
				IsUT:     ut,
			})
		}
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/states/all", geodata.ListStatesAndUT(false))
	mux.HandleFunc("/states/with-ut", geodata.ListStatesAndUT(true))
	mux.HandleFunc("/get-state", geodata.GetState)
	mux.HandleFunc("/states/east-west", geodata.OrderStates("EW"))
	mux.HandleFunc("/states/north-south", geodata.OrderStates("NS"))
	fmt.Println("Serving on :9000")
	log.Fatal(http.ListenAndServe(":9000", mux))

}

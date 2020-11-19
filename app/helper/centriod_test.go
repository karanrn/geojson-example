package helper

import (
	"testing"

	"github.com/karanrn/geojson-example/app/geodata"
	"github.com/paulmach/orb"
)

func TestGetCentriod(t *testing.T) {
	var polygon = orb.Ring{{0, 1}, {1, 1}, {1, 0}, {0, 0}}
	var expectedCentriod = geodata.Loc{Lat: 0.5, Lon: 0.5}
	result := GetCentriod(polygon)
	if result != expectedCentriod {
		t.Log("error should be (0,0), got", result)
		t.Fail()
	}
}

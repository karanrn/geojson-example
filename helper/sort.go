package helper

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

func (s StateCentriodList) Len() int {
	return len(s)
}

func (s StateCentriodList) Less(i, j int) bool {
	return s[i].Centriod.Lat > s[j].Centriod.Lat
}

func (s StateCentriodList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

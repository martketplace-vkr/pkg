package location

import (
	"time"

	"github.com/goccy/go-json"
)

type Location struct {
	*time.Location
}

func (l *Location) UnmarshalJSON(data []byte) error {
	var locationName string
	if err := json.Unmarshal(data, &locationName); err != nil {
		return err
	}
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return err
	}
	l.Location = loc
	return nil
}

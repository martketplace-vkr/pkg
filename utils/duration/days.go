package duration

import (
	"time"

	"github.com/goccy/go-json"
)

type Days struct {
	time.Duration
}

func (d *Days) MarshalJSON() ([]byte, error) {
	days := d.Duration / (SecondsInDay * time.Second)
	return json.Marshal(days)
}

func (d *Days) UnmarshalJSON(b []byte) error {
	var days int64
	if err := json.Unmarshal(b, &days); err != nil {
		return err
	}

	d.Duration = time.Duration(days) * SecondsInDay * time.Second

	return nil
}

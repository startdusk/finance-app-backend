package utils

import (
	"net/url"
	"time"
)

func TimeParam(query url.Values, name string) (time.Time, error) {
	// return NOW is default time.
	t := time.Now()
	value := query.Get(name)
	if value == "" {
		return t, nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return t, err
	}

	return parsed, nil
}

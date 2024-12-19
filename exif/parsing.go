package exif

import (
	"strconv"
	"strings"
)

func (few *exifFileExifWalker) captureExifGpsData(latLongSlice []string, direction string) float64 {
	var (
		degrees float64
		minutes float64
		seconds float64

		split []string
	)

	f := func(set string) float64 {
		split = strings.Split(set, "/")

		if len(split) < 2 {
			return 0.0
		}

		top, _ := strconv.ParseFloat(split[0], 64)
		bottom, _ := strconv.ParseFloat(split[1], 64)

		return top / bottom
	}

	if len(latLongSlice) == 3 && direction != "" {
		degrees = f(latLongSlice[0])
		minutes = f(latLongSlice[1])
		seconds = f(latLongSlice[2])

		calc := degrees + (minutes / 60) + (seconds / 3600)

		if direction == "S" || direction == "W" {
			calc *= -1
		}

		return calc
	}

	return 0.0
}

func stripOutsideQuotes(s string) string {
	result := strings.TrimPrefix(s, `"`)
	result = strings.TrimSuffix(result, `"`)

	return result
}

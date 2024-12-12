package imagemetadata

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type exifFileExifWalker struct {
	Image *ImageData

	latitude           []string
	latitudeDirection  string
	longitude          []string
	longitudeDirection string
}

func getExifData(input io.Reader, image *ImageData) error {
	x, err := exif.Decode(input)

	if err != nil {
		return fmt.Errorf("error reading EXIF data in file: %w", err)
	}

	walker := &exifFileExifWalker{Image: image}
	x.Walk(walker)
	return nil
}

func (few *exifFileExifWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	var (
		err error
	)

	switch name {
	case "Make":
		few.Image.Make = stripOutsideQuotes(tag.String())

	case "Model":
		few.Image.Model = stripOutsideQuotes(tag.String())

	case "LensModel":
		few.Image.LensModel = stripOutsideQuotes(tag.String())

	case "PixelXDimension":
		few.Image.Width, _ = strconv.Atoi(tag.String())

	case "PixelYDimension":
		few.Image.Height, _ = strconv.Atoi(tag.String())

	case "ImageDescription":
		few.Image.Caption = stripOutsideQuotes(tag.String())

	case "DateTimeOriginal":
		temp := stripOutsideQuotes(strings.TrimSpace(tag.String()))

		if temp != "" {
			split1 := strings.Split(temp, " ")

			if len(split1) == 2 {
				dateSplit := strings.Split(split1[0], ":")
				time := split1[1]

				few.Image.CreationDateTime = fmt.Sprintf("%s-%s-%sT%s", dateSplit[0], dateSplit[1], dateSplit[2], time)
			}
		}

	case "GPSLatitude":
		temp := stripOutsideQuotes(strings.TrimSpace(tag.String()))

		if temp != "" {
			if err = json.Unmarshal([]byte(temp), &few.latitude); err == nil {
				few.Image.Latitude = few.captureExifGpsData(few.latitude, few.latitudeDirection)
			}

		}

	case "GPSLatitudeRef":
		few.latitudeDirection = stripOutsideQuotes(strings.TrimSpace(tag.String()))
		few.Image.Latitude = few.captureExifGpsData(few.latitude, few.latitudeDirection)

	case "GPSLongitude":
		temp := stripOutsideQuotes(strings.TrimSpace(tag.String()))

		if temp != "" {
			if err = json.Unmarshal([]byte(temp), &few.longitude); err == nil {
				few.Image.Longitude = few.captureExifGpsData(few.longitude, few.longitudeDirection)
			}

		}

	case "GPSLongitudeRef":
		few.longitudeDirection = stripOutsideQuotes(strings.TrimSpace(tag.String()))
		few.Image.Longitude = few.captureExifGpsData(few.longitude, few.longitudeDirection)

	}

	return nil
}

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

package exif

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/adampresley/imagemetadata/imagemodel"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type exifFileExifWalker struct {
	Image   *imagemodel.ImageData
	GotEXIF bool

	latitude           []string
	latitudeDirection  string
	longitude          []string
	longitudeDirection string
}

func DecorateEXIF(image *imagemodel.ImageData, input io.ReadSeeker) (bool, error) {
	var (
		err               error
		decodedExifHeader *exif.Exif
	)

	if _, err = input.Seek(0, 0); err != nil {
		return false, fmt.Errorf("error resetting reader: %w", err)
	}

	if decodedExifHeader, err = exif.Decode(input); err != nil {
		return false, fmt.Errorf("error reading EXIF data in file: %w", err)
	}

	walker := &exifFileExifWalker{Image: image, GotEXIF: false}
	decodedExifHeader.Walk(walker)
	return walker.GotEXIF, nil
}

func (few *exifFileExifWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	var (
		err  error
		temp string
	)

	temp = stripOutsideQuotes(strings.TrimSpace(tag.String()))

	if temp != "" {
		few.GotEXIF = true
	}

	switch name {
	case "Make":
		if few.Image.Make == "" {
			few.Image.Make = temp
		}

	case "Model":
		if few.Image.Model == "" {
			few.Image.Model = temp
		}

	case "LensMake":
		if few.Image.LensMake == "" {
			few.Image.LensMake = temp
		}

	case "LensModel":
		if few.Image.LensModel == "" {
			few.Image.LensModel = temp
		}

	case "PixelXDimension":
		if few.Image.Width == 0 {
			few.Image.Width, _ = strconv.Atoi(temp)
		}

	case "PixelYDimension":
		if few.Image.Height == 0 {
			few.Image.Height, _ = strconv.Atoi(temp)
		}

	case "ImageDescription":
		few.Image.CaptionEXIF = temp

	case "DateTimeOriginal":
		if temp != "" && few.Image.CreationDateTime == "" {
			split1 := strings.Split(temp, " ")

			if len(split1) == 2 {
				dateSplit := strings.Split(split1[0], ":")
				time := split1[1]

				few.Image.CreationDateTime = fmt.Sprintf("%s-%s-%sT%s", dateSplit[0], dateSplit[1], dateSplit[2], time)
			}
		}

	case "GPSLatitude":
		if temp != "" {
			if err = json.Unmarshal([]byte(temp), &few.latitude); err == nil {
				few.Image.Latitude = few.captureExifGpsData(few.latitude, few.latitudeDirection)
			}
		}

	case "GPSLatitudeRef":
		few.latitudeDirection = stripOutsideQuotes(strings.TrimSpace(tag.String()))
		few.Image.Latitude = few.captureExifGpsData(few.latitude, few.latitudeDirection)

	case "GPSLongitude":
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

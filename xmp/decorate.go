package xmp

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/adampresley/adamgokit/slices"
	"github.com/adampresley/imagemetadata/imagemodel"
	"github.com/honza/go-xmp/xmp"
)

func DecorateXMP(image *imagemodel.ImageData, input io.ReadSeeker) (bool, error) {
	var (
		err       error
		xmpHeader []byte
	)

	if _, err = input.Seek(0, 0); err != nil {
		return false, fmt.Errorf("error resetting reader: %w", err)
	}

	if xmpHeader, err = getXMPDataBlock(input); err != nil {
		return false, fmt.Errorf("error reading XMP header: %w", err)
	}

	if xmpHeader == nil {
		return false, nil
	}

	decoder := xmp.NewDecoder(bytes.NewReader(xmpHeader))
	doc := &xmp.Document{}
	gotXMP := false

	if err = decoder.Decode(doc); err != nil {
		return false, fmt.Errorf("error decoding XMP data: %w", err)
	}

	xmpKeywords := xmpGetKeywords(doc)
	xmpPeople := xmpGetPeople(doc)
	xmpTitle := strings.TrimSpace(xmpGetTitle(doc))
	xmpLensModel := strings.TrimSpace(xmpGetLensModel(doc))
	xmpLat, xmpLong := xmpGetLatitudeLongitude(doc)
	lat := xmpGpsCoordToFloat(xmpLat)
	long := xmpGpsCoordToFloat(xmpLong)

	image.Keywords = slices.Merge(image.Keywords, xmpKeywords)
	image.People = slices.Merge(image.People, xmpPeople)

	if image.TitleXMP == "" {
		image.TitleXMP = xmpTitle
	}

	if image.LensModel == "" {
		image.LensModel = xmpLensModel
	}

	if image.Latitude == 0 {
		image.Latitude = lat
	}

	if image.Longitude == 0 {
		image.Longitude = long
	}

	gotXMP = len(xmpKeywords) > 0 || len(xmpPeople) > 0 ||
		xmpTitle != "" || xmpLensModel != "" || lat != 0.0 ||
		long != 0.0

	return gotXMP, nil
}

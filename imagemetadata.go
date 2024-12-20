package imagemetadata

import (
	"fmt"
	"io"
	"strings"

	"github.com/adampresley/adamgokit/slices"
	"github.com/adampresley/imagemetadata/exif"
	"github.com/adampresley/imagemetadata/imagemodel"
	"github.com/adampresley/imagemetadata/iptc"
	"github.com/adampresley/imagemetadata/xmp"
)

func NewFromJPEG(input io.ReadSeeker) (*imagemodel.ImageData, error) {
	var (
		err error

		gotIPTC bool
		gotEXIF bool
		gotXMP  bool
	)

	result := &imagemodel.ImageData{}
	errors := []error{}

	if gotIPTC, err = iptc.DecorateIPTC(result, input); err != nil {
		errors = append(errors, fmt.Errorf("error reading IPTC data: %w", err))
	}

	if gotEXIF, err = exif.DecorateEXIF(result, input); err != nil {
		errors = append(errors, fmt.Errorf("error reading EXIF data: %w", err))
	}

	if gotXMP, err = xmp.DecorateXMP(result, input); err != nil {
		errors = append(errors, fmt.Errorf("error reading XMP data: %w", err))
	}

	gotAnythingUseful := gotIPTC || gotEXIF || gotXMP

	if !gotAnythingUseful {
		errorMessages := strings.Join(slices.Map(errors, func(e error, index int) string { return e.Error() }), " :: ")
		err = fmt.Errorf("no metadata found: %s", errorMessages)
	}

	return result, err
}

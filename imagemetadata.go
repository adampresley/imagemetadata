package imagemetadata

import (
	"fmt"
	"io"

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

	if gotIPTC, err = iptc.DecorateIPTC(result, input); err != nil {
		return result, fmt.Errorf("error reading IPTC data: %w", err)
	}

	if gotEXIF, err = exif.DecorateEXIF(result, input); err != nil {
		return result, fmt.Errorf("error reading EXIF data: %w", err)
	}

	if gotXMP, err = xmp.DecorateXMP(result, input); err != nil {
		return result, fmt.Errorf("error reading XMP data: %w", err)
	}

	gotAnythingUseful := gotIPTC || gotEXIF || gotXMP

	if !gotAnythingUseful {
		err = fmt.Errorf("no metadata found")
	}

	return result, err
}

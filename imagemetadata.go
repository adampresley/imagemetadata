package imagemetadata

import (
	"fmt"
	"io"

	"github.com/adampresley/imagemetadata/exif"
	"github.com/adampresley/imagemetadata/imagemodel"
	"github.com/adampresley/imagemetadata/iptc"
	"github.com/adampresley/imagemetadata/xmp"
)

type ReadResult struct {
	Metadata *imagemodel.ImageData
	Errors   []error
}

func NewFromJPEG(input io.ReadSeeker) (ReadResult, error) {
	var (
		err      error
		finalErr error

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
		finalErr = fmt.Errorf("no metadata found")
	}

	return ReadResult{
		Metadata: result,
		Errors:   errors,
	}, finalErr
}

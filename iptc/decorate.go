package iptc

import (
	"bytes"
	"fmt"
	"io"

	"github.com/adampresley/imagemetadata/imagemodel"
)

func DecorateIPTC(imageData *imagemodel.ImageData, input io.ReadSeeker) (bool, error) {
	var (
		err       error
		dataBlock []byte
		tags      map[string][]string
	)

	if _, err = input.Seek(0, 0); err != nil {
		return false, fmt.Errorf("error resetting reader: %w", err)
	}

	if dataBlock, err = getIPTCDataBlock(input); err != nil {
		return false, fmt.Errorf("error getting data block: %w", err)
	}

	dataBlockReader := bytes.NewReader(dataBlock)

	if tags, err = readIPTC(dataBlockReader); err != nil {
		return false, fmt.Errorf("error parsing tag markers and data: %w", err)
	}

	for k, v := range tags {
		if f, ok := standardTags[k]; ok {
			f(imageData, v)
		}
	}

	return len(tags) > 0, nil
}

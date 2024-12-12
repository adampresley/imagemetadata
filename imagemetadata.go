package imagemetadata

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ImageData struct {
	Make             string
	Model            string
	LensMake         string
	LensModel        string
	Caption          string
	Title            string
	CreationDateTime string
	Width            int
	Height           int
	Latitude         float64
	Longitude        float64
	Keywords         []string
	People           []string
}

func NewFromJPEG(input io.ReadSeeker) (*ImageData, error) {
	var (
		err          error
		xmpDataBlock []byte
	)

	result := &ImageData{}

	if err = getExifData(input, result); err != nil {
		return result, err
	}

	if _, err = input.Seek(0, 0); err != nil {
		return result, fmt.Errorf("error resetting file reader to zero after reading EXIF data: %w", err)
	}

	if xmpDataBlock, err = getXMPDataBlock(input); err != nil {
		return result, fmt.Errorf("error reading XMP data block in file: %w", err)
	}

	if err = getXMPData(xmpDataBlock, result); err != nil {
		return result, fmt.Errorf("error getting XMP data: %w", err)
	}

	return result, nil
}

func (id ImageData) String() string {
	s := strings.Builder{}

	s.WriteString("Title: " + id.Title + "\n")
	s.WriteString("Make: " + id.Make + "\n")
	s.WriteString("Model: " + id.Model + "\n")
	s.WriteString("Lens Make: " + id.LensMake + "\n")
	s.WriteString("Lens Model: " + id.LensModel + "\n")
	s.WriteString("Caption: " + id.Caption + "\n")
	s.WriteString("Created Date: " + id.CreationDateTime + "\n")
	s.WriteString("Width: " + strconv.Itoa(id.Width) + "\n")
	s.WriteString("Height: " + strconv.Itoa(id.Height) + "\n")
	s.WriteString("Latitude: " + strconv.FormatFloat(id.Latitude, 'f', -1, 64) + "\n")
	s.WriteString("Longitude: " + strconv.FormatFloat(id.Longitude, 'f', -1, 64) + "\n")
	s.WriteString("Keywords: " + strings.Join(id.Keywords, ", ") + "\n")
	s.WriteString("People: " + strings.Join(id.People, ", ") + "\n")

	return s.String()
}

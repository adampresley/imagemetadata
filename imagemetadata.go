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
		err error
		// iptcDataBlock []byte
		xmpDataBlock []byte
	)

	result := &ImageData{}
	gotExif := true
	gotXMP := true

	// if iptcDataBlock, err = getIPTCDataBlock(input); err != nil {
	// 	return result, fmt.Errorf("error reading IPTC data block: %w", err)
	// }

	// fmt.Printf("IPTC Data Block: %x\n", iptcDataBlock)
	// if err = readIPTC(input); err != nil {
	// 	return result, fmt.Errorf("error reading IPTC data: %w", err)
	// }

	if _, err = input.Seek(0, 0); err != nil {
		return result, fmt.Errorf("error resetting file reader to zero after reading IPTC data: %w", err)
	}

	if err = getExifData(input, result); err != nil {
		gotExif = false
	}

	if _, err = input.Seek(0, 0); err != nil {
		return result, fmt.Errorf("error resetting file reader to zero after reading EXIF data: %w", err)
	}

	if xmpDataBlock, err = getXMPDataBlock(input); err != nil {
		gotXMP = false
	}

	if gotXMP {
		if err = getXMPData(xmpDataBlock, result); err != nil {
			gotXMP = false
		}
	}

	err = nil

	if !gotExif && !gotXMP {
		err = fmt.Errorf("no metadata found")
	}

	return result, err
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

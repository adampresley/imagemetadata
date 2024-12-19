package imagemodel

import (
	"strconv"
	"strings"
)

type ImageData struct {
	Address          string
	Author           string
	CaptionIPTC      string
	CaptionEXIF      string
	City             string
	Copyright        string
	CountryCode      string
	CountryName      string
	CreationDateTime string
	Headline         string
	Height           int
	Keywords         []string
	Latitude         float64
	LensMake         string
	LensModel        string
	Longitude        float64
	Make             string
	Model            string
	People           []string
	StateProvince    string
	TitleIPTC        string
	TitleXMP         string
	Width            int
}

func (id ImageData) String() string {
	s := strings.Builder{}

	s.WriteString("Address: " + id.Address + "\n")
	s.WriteString("Author: " + id.Author + "\n")
	s.WriteString("Caption (IPTC): " + id.CaptionIPTC + "\n")
	s.WriteString("Caption (EXIF): " + id.CaptionEXIF + "\n")
	s.WriteString("City: " + id.City + "\n")
	s.WriteString("Copyright: " + id.Copyright + "\n")
	s.WriteString("Country Code: " + id.CountryCode + "\n")
	s.WriteString("Country Name: " + id.CountryName + "\n")
	s.WriteString("Created Date: " + id.CreationDateTime + "\n")
	s.WriteString("Headline: " + id.Headline + "\n")
	s.WriteString("Height: " + strconv.Itoa(id.Height) + "\n")
	s.WriteString("Keywords: " + strings.Join(id.Keywords, ", ") + "\n")
	s.WriteString("Latitude: " + strconv.FormatFloat(id.Latitude, 'f', -1, 64) + "\n")
	s.WriteString("Lens Make: " + id.LensMake + "\n")
	s.WriteString("Lens Model: " + id.LensModel + "\n")
	s.WriteString("Longitude: " + strconv.FormatFloat(id.Longitude, 'f', -1, 64) + "\n")
	s.WriteString("Make: " + id.Make + "\n")
	s.WriteString("Model: " + id.Model + "\n")
	s.WriteString("People: " + strings.Join(id.People, ", ") + "\n")
	s.WriteString("StateProvince: " + id.StateProvince + "\n")
	s.WriteString("Title (IPTC): " + id.TitleIPTC + "\n")
	s.WriteString("Title (XMP): " + id.TitleXMP + "\n")
	s.WriteString("Width: " + strconv.Itoa(id.Width) + "\n")

	return s.String()
}

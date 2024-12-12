package imagemetadata

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/honza/go-xmp/xmp"
)

func getXMPDataBlock(input io.Reader) ([]byte, error) {
	var (
		err error
		n   int
	)

	startMarker := []byte(`<x:xmpmeta`)
	endMarker := []byte(`</x:xmpmeta>`)
	foundStart := false

	var xmpBuffer bytes.Buffer
	buffer := make([]byte, 1024)

	reader := bufio.NewReader(input)

	for {
		n, err = reader.Read(buffer)

		if err != nil && err != io.EOF {
			return []byte{}, fmt.Errorf("error performing read operation: %w", err)
		}

		if n == 0 {
			break
		}

		chunk := buffer[:n]

		if !foundStart {
			startIndex := bytes.Index(chunk, startMarker)

			if startIndex != -1 {
				foundStart = true
				xmpBuffer.Write(chunk[startIndex:])
			}
		} else {
			xmpBuffer.Write(chunk)
			endIdx := bytes.Index(chunk, endMarker)

			if endIdx != -1 {
				xmpBuffer.Truncate(xmpBuffer.Len() - len(chunk[endIdx+len(endMarker):]))
				break
			}
		}
	}

	if xmpBuffer.Len() == 0 {
		return []byte{}, fmt.Errorf("no xmp data found")
	}

	return xmpBuffer.Bytes(), nil
}

func getXMPData(xmpData []byte, image *ImageData) error {
	var (
		err error
	)

	decoder := xmp.NewDecoder(bytes.NewReader(xmpData))
	doc := &xmp.Document{}

	if err = decoder.Decode(doc); err != nil {
		return fmt.Errorf("error decoding XMP data: %w", err)
	}

	image.Keywords = xmpGetKeywords(doc)
	image.People = xmpGetPeople(doc)
	image.Title = xmpGetTitle(doc)

	if image.LensModel == "" {
		image.LensModel = xmpGetLensModel(doc)
	}

	lat, long := xmpGetLatitudeLongitude(doc)

	if image.Latitude == 0 {
		image.Latitude = xmpGpsCoordToFloat(lat)
	}

	if image.Longitude == 0 {
		image.Longitude = xmpGpsCoordToFloat(long)
	}

	return nil
}

func xmpGetKeywords(doc *xmp.Document) []string {
	result := []string{}
	ns := doc.FindNs("dc", "")
	node := doc.FindNode(ns)

	/*
	   * Keywords are found in <dc:subject><rdf:Bag><rdf:li>value</rdf:li></rdf:Bag></dc:subject>
	   * Example:
	   *   <rdf:Description rdf:about=''
	           xmlns:dc='http://purl.org/dc/elements/1.1/'>
	           <dc:subject>
	               <rdf:Bag>
	                   <rdf:li>2024</rdf:li>
	                   <rdf:li>Holiday</rdf:li>
	                   <rdf:li>Thanksgiving</rdf:li>
	               </rdf:Bag>
	           </dc:subject>
	       </rdf:Description>
	*/
	for _, container := range node.Nodes {
		if !container.IsZero() && container.XMLName.Local == "dc:subject" {
			for _, bag := range container.Nodes {
				for _, li := range bag.Nodes {
					result = append(result, li.Value)
				}
			}
		}
	}

	return result
}

func xmpGetPeople(doc *xmp.Document) []string {
	result := []string{}
	ns := doc.FindNs("Iptc4xmpExt", "")
	node := doc.FindNode(ns)

	if node != nil {
		for _, container := range node.Nodes {
			if !container.IsZero() && container.XMLName.Local == "Iptc4xmpExt:PersonInImage" {
				for _, bag := range container.Nodes {
					for _, li := range bag.Nodes {
						result = append(result, li.Value)
					}
				}
			}
		}
	}

	return result
}

func xmpGetLensModel(doc *xmp.Document) string {
	result := ""
	ns := doc.FindNs("aux", "")
	node := doc.FindNode(ns)

	if node != nil {
		if value, err := node.GetPath("aux:Lens"); err == nil {
			result = value
		}
	}

	return result
}

func xmpGetTitle(doc *xmp.Document) string {
	result := ""
	ns := doc.FindNs("dc", "")
	node := doc.FindNode(ns)

	if node != nil {
		if value, err := node.GetPath("dc:title/rdf:Alt/rdf:li"); err == nil {
			result = value
		}
	}

	return result
}

func xmpGetLatitudeLongitude(doc *xmp.Document) (string, string) {
	var (
		err              error
		lat, long, value string
	)
	ns := doc.FindNs("exif", "")

	if ns != nil {
		node := doc.FindNode(ns)

		if node != nil {
			if value, err = node.GetPath("exif:GPSLatitude"); err == nil {
				lat = value
			}

			if value, err = node.GetPath("exif:GPSLongitude"); err == nil {
				long = value
			}
		}
	}

	return lat, long
}

func xmpGpsCoordToFloat(coord string) float64 {
	var (
		err     error
		degrees float64
		minutes float64
	)

	if coord != "" {
		coord = strings.ReplaceAll(coord, ",", ".")

		split := strings.SplitN(coord, ".", 1)

		if len(split) == 2 {
			direction := split[1][len(split[1])-1:]
			minuteString := strings.TrimRightFunc(split[1], func(r rune) bool {
				if r == 'N' || r == 'W' || r == 'E' || r == 'S' {
					return true
				}

				return false
			})

			if degrees, err = strconv.ParseFloat(split[0], 64); err != nil {
				fmt.Printf("error converting '%s' to float\n", split[0])
				return 0.0
			}

			if minutes, err = strconv.ParseFloat(minuteString, 64); err != nil {
				fmt.Printf("error converting '%s' to float\n", minuteString)
				return 0.0
			}

			decimalDegrees := degrees + (minutes / 60)

			if direction == "S" || direction == "W" {
				decimalDegrees *= -1
			}

			return decimalDegrees
		}
	}

	return 0.0
}

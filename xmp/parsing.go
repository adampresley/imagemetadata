package xmp

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
		err        error
		n          int
		foundStart bool
		totalBytes int
		xmpBuffer  bytes.Buffer
	)

	startMarker := []byte(`<x:xmpmeta`)
	endMarker := []byte(`</x:xmpmeta>`)
	buffer := make([]byte, 1024)
	startIndex := -1
	endIndex := -1

	reader := bufio.NewReader(input)

	for {
		n, err = reader.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		if n == 0 { // End of file
			break
		}

		totalBytes += n
		if _, err = xmpBuffer.Write(buffer); err != nil {
			return nil, fmt.Errorf("error writing XMP chunk to buffer: %w", err)
		}

		chunk := xmpBuffer.Bytes()

		if !foundStart {
			// Look for the start marker
			startIndex = bytes.Index(chunk, startMarker)

			if startIndex != -1 {
				foundStart = true
			}
		} else {
			// Look for the end marker
			endIndex = bytes.Index(chunk, endMarker)

			if endIndex != -1 {
				break
			}
		}
	}

	if startIndex == -1 && endIndex == -1 {
		return nil, nil
		// return nil, fmt.Errorf("no XMP data found")
	}

	b := xmpBuffer.Bytes()
	header := b[startIndex : endIndex+len(endMarker)]
	return header, nil
}

func xmpGetKeywords(doc *xmp.Document) []string {
	result := []string{}
	ns := doc.FindNs("dc", "")

	if ns != nil {
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
		if node != nil {
			for _, container := range node.Nodes {
				if !container.IsZero() && container.XMLName.Local == "dc:subject" {
					for _, bag := range container.Nodes {
						for _, li := range bag.Nodes {
							result = append(result, li.Value)
						}
					}
				}
			}
		}
	}

	return result
}

func xmpGetPeople(doc *xmp.Document) []string {
	result := []string{}
	ns := doc.FindNs("Iptc4xmpExt", "")

	if ns != nil {
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
	}

	return result
}

func xmpGetLensModel(doc *xmp.Document) string {
	result := ""
	ns := doc.FindNs("aux", "")

	if ns != nil {
		node := doc.FindNode(ns)

		if node != nil {
			if value, err := node.GetPath("aux:Lens"); err == nil {
				result = value
			}
		}
	}

	return result
}

func xmpGetTitle(doc *xmp.Document) string {
	result := ""
	ns := doc.FindNs("dc", "")

	if ns != nil {
		node := doc.FindNode(ns)

		if node != nil {
			if value, err := node.GetPath("dc:title/rdf:Alt/rdf:li"); err == nil {
				result = value
			}
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

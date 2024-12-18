package imagemetadata

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dsoprea/go-iptc"
)

func readIPTC(input io.Reader) error {
	var (
		err        error
		data       map[string]string
		streamData map[iptc.StreamTagKey][]iptc.TagData
	)

	if streamData, err = iptc.ParseStream(input); err != nil {
		return fmt.Errorf("could not read IPTC stream data: %w", err)
	}

	data = iptc.GetSimpleDictionaryFromParsedTags(streamData)
	fmt.Printf("IPTC: %+v\n", data)

	return nil
}

func getIPTCDataBlock(input io.Reader) ([]byte, error) {
	const bufferSize = 1024
	buffer := make([]byte, bufferSize)

	var iptcData bytes.Buffer
	foundIPTCStart := false

	for {
		// Read a chunk of the file
		n, err := input.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		chunk := buffer[:n]
		if !foundIPTCStart {
			// Search for the IPTC start marker (0x1C)
			iptcStartIndex := bytes.Index(chunk, []byte{0x1C})
			if iptcStartIndex != -1 {
				foundIPTCStart = true
				iptcData.Write(chunk[iptcStartIndex:])
			}
		} else {
			// Continue collecting data after the start marker
			iptcData.Write(chunk)
		}
	}

	if !foundIPTCStart {
		return nil, fmt.Errorf("IPTC data not found")
	}

	return iptcData.Bytes(), nil
}

package iptc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func readIPTC(input *bytes.Reader) (map[string][]string, error) {
	tags := make(map[string][]string)

	for input.Len() > 0 {
		// Read tag marker (should be 0x1C). If not, skip.
		marker, err := input.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed to read IPTC tag marker: %w", err)
		}

		if marker != 0x1C {
			continue
		}

		// Read dataset and record numbers
		dataset, err := input.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed to read dataset: %w", err)
		}
		record, err := input.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		// Read length of the data
		var length uint16
		if err := binary.Read(input, binary.BigEndian, &length); err != nil {
			return nil, fmt.Errorf("failed to read data length: %w", err)
		}

		// Read the data value
		value := make([]byte, length)
		if _, err := input.Read(value); err != nil {
			return nil, fmt.Errorf("failed to read tag value: %w", err)
		}

		// Construct a key for the tag (e.g., "2:25" for dataset 2, record 25)
		tagKey := fmt.Sprintf("%d:%d", dataset, record)

		if _, ok := tags[tagKey]; !ok {
			tags[tagKey] = []string{}
		}

		tags[tagKey] = append(tags[tagKey], string(value))
	}

	return tags, nil
}

func getIPTCDataBlock(input io.ReadSeeker) ([]byte, error) {
	var (
		err             error
		i               int64
		blockLength     int64 = 0
		blockStartIndex int64 = 0
	)

	/*
	 * Read the header till we find the 0xFFED marker. From here we can
	 * get the length of the IPTC header. Use that to read the entire block.
	 */
	buffer := make([]byte, 1024)

	for {
		_, err := input.Read(buffer)
		if err != nil {
			break
		}

		// Search for APP13 marker (0xFFED)
		for i = 0; i < int64(len(buffer)-1); i++ {
			if buffer[i] == 0xFF && buffer[i+1] == 0xED && (i+4) < int64(len(buffer)) {
				// Length of APP13 segment (next 2 bytes, big-endian)
				blockLength = int64(binary.BigEndian.Uint16(buffer[i+2 : i+4]))
				blockStartIndex += i + 4
				break
			}
		}

		if blockStartIndex > 0 && blockLength > 0 {
			break
		}

		blockStartIndex += int64(len(buffer))
	}

	if blockStartIndex > 0 && blockLength > 0 {
		blockBuffer := make([]byte, blockLength)

		if _, err = input.Seek(blockStartIndex, 0); err != nil {
			return nil, fmt.Errorf("error resetting reader when getting block data: %w", err)
		}

		if _, err = input.Read(blockBuffer); err != nil {
			return nil, fmt.Errorf("error reading block buffer (len %d, offset %d): %w", blockLength, blockStartIndex, err)
		}

		firstMarker := bytes.Index(blockBuffer, []byte("8BIM"))
		lastMarker := bytes.LastIndex(blockBuffer, []byte("8BIM"))

		if lastMarker > firstMarker && (lastMarker+len([]byte("8BIM"))) < int(blockLength) {
			return blockBuffer[:lastMarker], nil
		}

		return blockBuffer, nil
	}

	return nil, fmt.Errorf("IPTC data not found in file")
}

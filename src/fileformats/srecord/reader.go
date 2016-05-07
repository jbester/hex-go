package srecord

import (
	"io"
	"fmt"
	"bufio"
	"strings"
	"encoding/binary"
)

type SRecordReader struct {
	reader *bufio.Reader
	lineNo int
}


func Open(reader io.Reader) *SRecordReader {
	result := SRecordReader{
		reader: bufio.NewReader(reader),
	}
	return &result
}

func (self *SRecordReader) readLine() (string, error) {
	s, err := self.reader.ReadString('\n')
	if err != nil {
		return s, err
	}
	self.lineNo += 1
	return strings.TrimRight(s, "\r\n"), nil
}



func (self *SRecordReader) parseDataRecord(addressWidth int, record string) (*DataRecord, error) {
	addressBytes := addressWidth/8
	recordAsBytes, err := hexToBytes(record[2:])
	if err != nil {
		return nil, err
	}
	recordCount := recordAsBytes[0]
	address := recordAsBytes[1:1+addressBytes]
	data := recordAsBytes[1+addressBytes:len(recordAsBytes)-1]
	checksum := recordAsBytes[len(recordAsBytes)-1]
	actualChecksum := byte(0)
	actualChecksum += recordCount
	for _, b := range(address) {
		actualChecksum += b
	}

	for _, b := range(data) {
		actualChecksum += b
	}

	actualChecksum = actualChecksum ^ 0xFF
	if checksum != actualChecksum {
		return nil, fmt.Errorf("Checksum mismatch on line %v expected: %02X got: %02X",
			self.lineNo, checksum, actualChecksum)
	}

	for len(address) < 4 {
		address = append([]byte{0}, address...)
	}

	result := DataRecord{
		data: data,
		address: binary.BigEndian.Uint32(address),
	}
	return &result, nil
}

func (self *SRecordReader) Next() (*DataRecord, error) {
	for {
		line, err := self.readLine()
		if err != nil {
			return nil, err
		}

		// only parse data records
		switch {
		case strings.HasPrefix(line, "S1"):
			return self.parseDataRecord(16, line)

		case strings.HasPrefix(line, "S2"):
			return self.parseDataRecord(24, line)

		case strings.HasPrefix(line, "S3"):
			return self.parseDataRecord(32, line)
		}
	}
}

type IterationResult struct {
	Record *DataRecord
	Error error
}

func (reader *SRecordReader) Records() <-chan IterationResult {
	ch := make(chan IterationResult)
	go func () {
		defer close(ch)
		for {
			record, err := reader.Next()
			if err == io.EOF {
				break
			}
			ch<- IterationResult{
				Record: record,
				Error: err,
			}
			if err != nil {
				break
			}
		}
	}()
	return ch
}


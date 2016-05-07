/*
Copyright (c) 2016,  Jeffrey Bester

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package srecord

import (
	"io"
	"fmt"
	"bufio"
	"strings"
	"encoding/binary"
)

// S-Record reader
type SRecordReader struct {
	reader *bufio.Reader
	lineNo int
}

// create an S-Record reader form the provided reader
func Open(reader io.Reader) *SRecordReader {
	result := SRecordReader{
		reader: bufio.NewReader(reader),
	}
	return &result
}

// read and trim a line from the file
func (self *SRecordReader) readLine() (string, error) {
	s, err := self.reader.ReadString('\n')
	if err != nil {
		return s, err
	}
	self.lineNo += 1
	return strings.TrimRight(s, "\r\n"), nil
}



// parse the file
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

// Retrieve the next data record from the reader
// returns nil, io.EOF on end of file
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

// Iteration record - contains data or error
type IterationResult struct {
	Record *DataRecord
	Error error
}

// Records provides iteration over all records
// returns channel providing IterationResult
// each with either the Record or Error set
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


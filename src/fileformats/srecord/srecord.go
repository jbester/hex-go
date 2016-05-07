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

type RecordType uint8

const (
	S0 RecordType = iota
	S1
	S2
	S3
	S4
	S5
	S6
	S7
	S8
	S9
	Header = S0
	Data16BitAddress = S1
	Data24BitAddress = S2
	Data32BitAddress = S3
	RecordCount16Bit = S5
	RecordCount24Bit = S6
	StartAddress32Bit = S7
	StartAddress24Bit= S8
	StartAddress16Bit = S9
)

type HeaderRecord struct {
	data []byte
}

func (record *HeaderRecord) Data() []byte {
	return record.data
}

type DataRecord struct {
	address uint32
	data []byte
}

func (record *DataRecord) Data() []byte {
	return record.data
}

func (record *DataRecord) Address() uint32 {
	return record.address
}


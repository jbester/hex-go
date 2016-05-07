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
	"testing"
	"bytes"
	"io"
)

// Verify address and data match expected value
func TestReadOneRecord(t *testing.T) {
	buffer := bytes.NewBufferString(s19TestFile)
	reader := Open(buffer)
	if reader != nil {
		rec, err := reader.Next()
		if err != nil {
			t.Fatal(err)
		}
		if rec == nil {
			t.Fatalf("Returned record is nil")
		}
		if rec.Address() != 0x400 {
			t.Errorf("Address mismatch expected: 0x400 got: %v", rec.Address())
		}
		data := rec.Data()
		if bytes.Compare(data, rec.Data()[:len(data)]) != 0 {
			t.Errorf("Data mismatch")
		}
	} else {
		t.Fatal("Open call failed")
	}
}

// Verify address and data are expected for entire file
func TestNext(t *testing.T) {
	buffer := bytes.NewBufferString(s19TestFile)
	reader := Open(buffer)
	if reader != nil {
		data := make([]byte, 0)
		expectedAddress := uint32(0x400)
		for {
			rec, err := reader.Next()
			if err == io.EOF {
				break
			}
			if rec.Address() != expectedAddress {
				t.Fatalf("Address mismatch expected: %v got: %v", expectedAddress, rec.Address())
			}
			expectedAddress += uint32(len(rec.Data()))
			data = append(data, rec.Data()...)
		}
		if bytes.Compare(data, binaryTestFile[:]) != 0 {
			t.Error("Data read did not match reference values")
		}
	} else {
		t.Fatal("Open call failed")
	}
}


// Verify iteration generates entire file
func TestIterate(t *testing.T) {
	buffer := bytes.NewBufferString(s19TestFile)
	reader := Open(buffer)
	if reader != nil {
		data := make([]byte, 0)
		expectedAddress := uint32(0x400)
		for it := range reader.Records() {
			rec := it.Record
			err := it.Error
			if err == io.EOF {
				t.Fatal("EOF not handled properly")
			}
			if rec.Address() != expectedAddress {
				t.Fatalf("Address mismatch expected: %v got: %v", expectedAddress, rec.Address())
			}
			expectedAddress += uint32(len(rec.Data()))
			data = append(data, rec.Data()...)
		}
		if bytes.Compare(data, binaryTestFile[:]) != 0 {
			t.Error("Data read did not match reference values")
		}
	} else {
		t.Fatal("Open call failed")
	}
}

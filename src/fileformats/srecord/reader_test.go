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

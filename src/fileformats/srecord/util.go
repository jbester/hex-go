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

import "strconv"

// Parse from hex to a byte
func hexToByte(hex string) (byte, error) {
	value, err := strconv.ParseUint(hex, 16, 8)
	return byte(value), err

}


// Parse from multiple hex nibbles to a byte slice
func hexToBytes(hex string) ([]byte, error) {
	result := make([]byte, 0)
	for i := 0; i < len(hex); i+=2 {
		value, err := hexToByte(hex[i:i+2])
		if err != nil {
			return nil, err
		}
		result = append(result, byte(value))
	}

	return result, nil

}

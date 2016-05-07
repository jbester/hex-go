package srecord

import "strconv"

func hexToByte(hex string) (byte, error) {
	value, err := strconv.ParseUint(hex, 16, 8)
	return byte(value), err

}


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

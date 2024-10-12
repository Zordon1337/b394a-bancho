package Packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func SerializePacket(packet int16, data []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, packet)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, bool(false))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, int32(len(data)))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, (data))
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// https://github.com/lenforiee/gocho/blob/main/packets/writer.go
func WriteUleb128(value int) (ret []byte) {
	var len int

	if value == 0 {
		ret = []byte{0}
		return ret
	}

	for value > 0 {
		ret = append(ret, 0)
		ret[len] = byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			ret[len] |= 0x80
		}
		len++
	}
	return ret
}

// https://github.com/lenforiee/gocho/blob/main/packets/writer.go
func WriteOsuString(value string) (ret []byte) {

	if value == "" {
		ret = []byte{0}
	} else {
		b := []byte(value)
		ret = append(ret, 0x0B)
		ret = append(ret, WriteUleb128(len(b))...)
		ret = append(ret, b...)
	}

	return ret
}
func ReadUleb128(r *bytes.Reader) (int, error) {
	var result int
	var shift uint
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		result |= int(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return result, nil
}

func ReadOsuString(r *bytes.Reader) (string, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return "", err
	}

	if prefix == 0x00 {
		// Empty string
		return "", nil
	} else if prefix == 0x0B {
		strLen, err := ReadUleb128(r)
		if err != nil {
			return "", err
		}

		// Read the string data of length strLen
		strData := make([]byte, strLen)
		if _, err := io.ReadFull(r, strData); err != nil {
			return "", err
		}

		return string(strData), nil
	}

	return "", fmt.Errorf("invalid osu! string prefix: %x", prefix)
}

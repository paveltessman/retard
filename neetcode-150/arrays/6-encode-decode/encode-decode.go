package encodedecode

import (
	"encoding/binary"
)

func encode(strs []string) string {
	buf := make([]byte, 0)

	for _, word := range strs {
		buf = binary.BigEndian.AppendUint32(buf, uint32(len(word)))
		buf = append(buf, []byte(word)...)
	}

	return string(buf)
}

func decode(s string) []string {
	encoded := []byte(s)
	result := make([]string, 0)

	var i uint32
	for i < uint32(len(encoded)) {

		length := binary.BigEndian.Uint32(encoded[i : i+4])
		i += 4

		word := encoded[i : i+length]
		result = append(result, string(word))
		i = i + length

	}
	return result
}

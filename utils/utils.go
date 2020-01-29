package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

// HandleErr : Handle error function
func HandleErr(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

// ToHex function
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

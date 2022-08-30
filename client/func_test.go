package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	//var pkgLen uint32
	//pkgLen = uint32(81)
	//fmt.Println("pkgLen=", pkgLen)
	//var buf [4]byte
	//fmt.Printf("buf=%v,buf[:4]=%v\n", buf, buf[:4])
	//binary.BigEndian.PutUint32(buf[0:4], pkgLen) //
	//
	//fmt.Println("buf", buf)

	buf := [4]byte{0, 2, 0, 81}

	x1 := binary.BigEndian.Uint32(buf[:4])
	fmt.Println("x1", x1, buf[:3])

}

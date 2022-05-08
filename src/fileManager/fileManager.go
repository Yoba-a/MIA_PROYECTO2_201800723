package fileManager

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"structs"
	"unsafe"
)

func Getfree(superBloque structs.Superbloque, path string, t string) int {
	var ch int8 = 0
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if t == "BI" {
		file.Seek(superBloque.S_bm_inode_start, 0)
		for k := 0; k < int(superBloque.S_inodes_count); k++ {
			data := leerBytes(file, int64(unsafe.Sizeof(ch)))
			buffer := bytes.NewBuffer(data)
			err = binary.Read(buffer, binary.BigEndian, &ch)
			if err != nil {
				log.Fatal("leer archivobinary.Read failed", err)
			}
			if ch == 0 {
				file.Close()
				return k
			}
		}
	} else {
		file.Seek(superBloque.S_bm_block_start, 0)
		for k := 0; k < int(superBloque.S_blocks_count); k++ {
			data := leerBytes(file, int64(unsafe.Sizeof(ch)))
			buffer := bytes.NewBuffer(data)
			err = binary.Read(buffer, binary.BigEndian, &ch)
			if err != nil {
				log.Fatal("leer archivobinary.Read failed", err)
			}
			if ch == 0 {
				file.Close()
				return k
			}
		}
	}
	file.Close()
	return -1
}

func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("ERROR  A LEER BYTES", err)
	}

	return bytes
}

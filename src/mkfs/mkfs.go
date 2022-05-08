package mkfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"structs"
	"time"
	"unsafe"
)

var path = ""
var name string = ""
var id = ""
var disco_amontar = structs.DiscoMontado{}
var tipo_ = ""

var type_com = regexp.MustCompile("(?i)\\s?-\\s?type\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var id_com = regexp.MustCompile("(?i)\\s?-\\s?id\\s?=\\s?([0-9]{3}[A-Z])")
var ruta_dsk = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")

var name_ = regexp.MustCompile("([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var name_id = regexp.MustCompile("([0-9]{3}[A-Z])")

var masterBoot = structs.Mbr{}
var Ebr = structs.Ebr{}

//variables para verificar la existencia del archivo
var pos2 = 0
var abs_path = ""

var path_disco = ""

func Analizador(input string) {
	if id_com.MatchString(input) && type_com.MatchString(input) {
		id = name_id.FindString(id_com.FindString(input))
		tipo_ = name_.FindString(type_com.FindString(input))

		fmt.Printf("name = %s\n", name)
		fmt.Printf("id = %s\n", id)
	} else {
		fmt.Println("error sintaxis no esperada los siguientes parametros son obligatorios: ")
		fmt.Println("valores no reconocidos -path: ")
		fmt.Printf("%q\n", type_com.Split(input, -1))
		fmt.Println("valores no reconocidos -name: ")
		fmt.Printf("%q\n", id_com.Split(input, -1))
	}

}

func FormatearExt2() {
	var dato string = string(id[2])
	disk, _ := strconv.Atoi(dato)
	var disco_nuevo = structs.DiscoMontado{}
	fmt.Println("Formateando la particion a EXT2")
	if structs.Discos_montados().Len() > 0 {
		for k := structs.Discos_montados().Front(); k != nil; k = k.Next() {
			disco_iterado := structs.DiscoMontado(k.Value.(structs.DiscoMontado))
			if disco_iterado.ID == disk {
				disco_nuevo = disco_iterado
				break
			}
		}
		path = ruta_dsk.FindString(string(disco_nuevo.Path[:]))
		var name_partition = ""

		for k := range disco_nuevo.Lista {
			var id_part_mount = string(disco_nuevo.Lista[k].ID[:])
			if id_part_mount == id {
				name_partition = string(disco_nuevo.Lista[k].Nombre[:])
				break
			}
		}
		Abrir_mbr()
		var iterador = -1
		var ext = -1
		for s := range masterBoot.Tabla {
			var name_par = string(masterBoot.Tabla[s].Name[:])
			if name_par == name_partition {
				iterador = s
			}
			if masterBoot.Tabla[s].Type == 'e' {
				ext = s
			}
		}
		if iterador == -1 && ext != -1 {
			/*
				floor((partition.part_size - sizeof(Structs::Superblock)) /
				(4 + sizeof(Structs::Inodes) + 3 * sizeof(Structs::Fileblock)))
			*/
			//cambiar para que pueda retornar datos de primarias o logicas
		} else {
			var part_size = masterBoot.Tabla[iterador].Size
			var equ = part_size - int64(unsafe.Sizeof(structs.Superbloque{}))
			var result = equ / (4 + int64(unsafe.Sizeof(structs.Inodes{})) + 3*int64(unsafe.Sizeof(structs.Fileblock{})))

			var n = int64(math.Floor(float64(result)))

			superBloque := structs.Superbloque{}
			superBloque.S_inodes_count = n
			superBloque.S_free_inodes_count = n
			superBloque.S_blocks_count = int64(3) * n
			superBloque.S_free_blocks_count = int64(3) * n
			currentTime := time.Now()
			copy(superBloque.Mtime[:], currentTime.Format("2006-01-02 15:04:05"))
			superBloque.S_mnt_count = 1
			superBloque.S_filesystem_type = 2

		}

	} else {
		fmt.Println("no se puede realizar el file system debido a que no hay ni un disco montado ")
	}
}

func ext2(superbloque structs.Superbloque, particion structs.Particion, n int64) {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	abs_path = path[:pos2]

	if ArchivoExiste(path) {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		superbloque.S_bm_inode_start = particion.Start + int64(unsafe.Sizeof(structs.Superbloque{}))
		superbloque.S_bm_block_start = superbloque.S_bm_inode_start + n
		superbloque.S_inode_start = superbloque.S_bm_block_start + (int64(3) * n)
		superbloque.S_bm_block_start = superbloque.S_bm_inode_start + (n + int64(unsafe.Sizeof(structs.Inodes{})))
		file.Seek(particion.Start, 0)
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, superbloque)
		writeNextBytes(file, binario3.Bytes())

		//formateando
		var vacio int8 = 0
		s := &vacio
		var binario bytes.Buffer
		file.Seek(superbloque.S_bm_inode_start, 0)
		binary.Write(&binario, binary.BigEndian, s)
		writeNextBytes(file, binario.Bytes())

		//situando el cursor en la ultima posicion
		file.Seek(n, 0)

		//colocando el ultimo byte para rellenar
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		writeNextBytes(file, binario2.Bytes())

		file.Seek(superbloque.S_bm_block_start, 0)
		binary.Write(&binario, binary.BigEndian, s)
		writeNextBytes(file, binario.Bytes())

		file.Seek(n*3, 0)
		binary.Write(&binario2, binary.BigEndian, s)
		writeNextBytes(file, binario2.Bytes())
		//fin de formateo

		//inodo
		file.Seek(superbloque.S_inode_start, 0)
		binary.Write(&binario3, binary.BigEndian, structs.Inodes{})
		writeNextBytes(file, binario3.Bytes())

		//folder
		file.Seek(superbloque.S_block_start, 0)
		binary.Write(&binario3, binary.BigEndian, structs.FolderBlock{})
		writeNextBytes(file, binario3.Bytes())

		file.Close()

	} else {
		fmt.Print("error el disco no existe en esta computadora...")
	}
}

func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func Abrir_ebr(inicio_ebr int64) {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	abs_path = path[:pos2]

	if ArchivoExiste(path) {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		var tamano_ebr int64 = int64(unsafe.Sizeof(Ebr))
		file.Seek(inicio_ebr, 0)
		data := leerBytes(file, tamano_ebr)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &Ebr)
		if err != nil {
			log.Fatal("leer archivobinary.Read failed", err)
		}
		fmt.Printf("nombre del ebr : %s \n", Ebr.Name[:])
		fmt.Printf("status ebr : %d \n", Ebr.Status)
		fmt.Printf("fit ebr : %s \n", Ebr.Fit)
		fmt.Printf("siguiente ebr : %d \n", Ebr.Next)
		fmt.Printf("size ebr : %d \n", Ebr.Size)
		fmt.Printf("inicio ebr : %d \n", Ebr.Start)
		file.Close()

	} else {
		fmt.Print("error el disco no existe en esta computadora...")
	}
}

func Abrir_mbr() {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	abs_path = path[:pos2]

	if ArchivoExiste(path) {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		var tamano_masterBoot int64 = int64(unsafe.Sizeof(masterBoot))
		data := leerBytes(file, tamano_masterBoot)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &masterBoot)
		if err != nil {
			log.Fatal("leer archivobinary.Read failed", err)
		}
		/*
			fmt.Println("Mbr Tamano:", masterBoot.Tamano)
			fmt.Println("Mbr Fecha creacion:", string(masterBoot.Fecha[:]))
			fmt.Println("Mbr Disk Signarue:", masterBoot.Firma)
			fmt.Println("Disco Fit:", string(masterBoot.Fit[:]))
			for k := range masterBoot.Tabla {
				fmt.Println("particion:", string(masterBoot.Tabla[k].Name[:]))
				fmt.Println("size :", masterBoot.Tabla[k].Size)
			}*/
		file.Close()

	} else {
		fmt.Print("error el disco no existe en esta computadora...")
	}
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("ERROR  A LEER BYTES", err)
	}

	return bytes
}

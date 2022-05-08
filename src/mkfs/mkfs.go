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
			fmt.Println("para particiones logicas aun nel xd")
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

			ext2(masterBoot.Tabla[iterador], n)
		}

	} else {
		fmt.Println("no se puede realizar el file system debido a que no hay ni un disco montado ")
	}
}

func ext2(particion structs.Particion, n int64) {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	abs_path = path[:pos2]

	if ArchivoExiste(path) {
		file, err := os.OpenFile(path, os.O_RDWR, 0)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		var p = int64(unsafe.Sizeof(masterBoot))
		fmt.Println(p)
		var superbloque = structs.Superbloque{}
		superbloque.S_inodes_count = n
		superbloque.S_free_inodes_count = n
		superbloque.S_blocks_count = int64(3) * n
		superbloque.S_free_blocks_count = int64(3) * n
		currentTime := time.Now()
		copy(superbloque.Mtime[:], currentTime.Format("2006-01-02 15:04:05"))
		superbloque.S_mnt_count = 1
		superbloque.S_filesystem_type = 2
		superbloque.S_bm_inode_start = particion.Start + int64(unsafe.Sizeof(structs.Superbloque{}))
		superbloque.S_bm_block_start = superbloque.S_bm_inode_start + n
		superbloque.S_inode_start = superbloque.S_bm_block_start + (int64(3) * n)
		superbloque.S_block_start = superbloque.S_bm_inode_start + (n + int64(unsafe.Sizeof(structs.Inodes{})))
		file.Seek(particion.Start, 0)
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, &superbloque)
		writeNextBytes(file, binario.Bytes())

		//formateando
		var vacio int8 = 0
		s := &vacio
		var binario2 bytes.Buffer
		file.Seek(superbloque.S_bm_inode_start, 0)
		binary.Write(&binario2, binary.BigEndian, s)
		writeNextBytes(file, binario2.Bytes())

		//situando el cursor en la ultima posicion
		file.Seek(n, 0)

		//colocando el ultimo byte para rellenar
		var binario_ bytes.Buffer
		binary.Write(&binario_, binary.BigEndian, s)
		writeNextBytes(file, binario_.Bytes())

		file.Seek(superbloque.S_bm_block_start, 0)
		var binario__ bytes.Buffer
		binary.Write(&binario__, binary.BigEndian, s)
		writeNextBytes(file, binario__.Bytes())

		file.Seek(n*3, 0)
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s)
		writeNextBytes(file, binario3.Bytes())
		//fin de formateo

		//inodo
		file.Seek(superbloque.S_inode_start, 0)
		for i := 0; i < int(n); i++ {
			var binario4 bytes.Buffer
			binary.Write(&binario4, binary.BigEndian, structs.Inodes{})
			writeNextBytes(file, binario4.Bytes())
		}

		//folder
		file.Seek(superbloque.S_block_start, 0)
		for j := 0; j < int(3*n); j++ {
			var binario5 bytes.Buffer
			binary.Write(&binario5, binary.BigEndian, structs.FolderBlock{})
			writeNextBytes(file, binario5.Bytes())
		}
		block_minus := [15]int64{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
		inode := structs.Inodes{}
		inode.I_uid = 1
		inode.I_gid = 1
		inode.I_size = 0
		inode.I_atime = superbloque.Mtime
		inode.I_ctime = superbloque.Mtime
		inode.I_mtime = superbloque.Mtime
		inode.I_type = 0
		inode.I_perm = 664
		inode.I_block = block_minus
		inode.I_block[0] = 0

		fb := structs.FolderBlock{}
		fb.B_content[0].B_name[0] = '.'
		copy(fb.B_content[1].B_name[:], "..")
		copy(fb.B_content[2].B_name[:], "user.txt")
		fb.B_content[2].B_inodo = 1

		var data = "1,G,root\n1,U, root,123\n"

		inodetmp := structs.Inodes{}
		inodetmp.I_uid = 1
		inodetmp.I_gid = 1
		inodetmp.I_size = int64(unsafe.Sizeof(data) + unsafe.Sizeof(structs.FolderBlock{}))
		inodetmp.I_atime = superbloque.Mtime
		inodetmp.I_ctime = superbloque.Mtime
		inodetmp.I_mtime = superbloque.Mtime
		inodetmp.I_type = 1
		inodetmp.I_perm = 664
		inodetmp.I_block = block_minus
		inodetmp.I_block[0] = 1

		inode.I_size = inodetmp.I_size + int64(unsafe.Sizeof(structs.FolderBlock{})) + int64(unsafe.Sizeof(structs.Inodes{}))

		fileb := structs.Fileblock{}
		copy(fileb.B_content[:], data)
		file.Seek(superbloque.S_bm_inode_start, 0)
		var char int8 = 1
		var binario6 bytes.Buffer
		binary.Write(&binario6, binary.BigEndian, char)
		writeNextBytes(file, binario6.Bytes())
		binary.Write(&binario6, binary.BigEndian, char)
		writeNextBytes(file, binario6.Bytes())

		file.Seek(superbloque.S_bm_block_start, 0)
		var binario7 bytes.Buffer
		binary.Write(&binario7, binary.BigEndian, char)
		writeNextBytes(file, binario7.Bytes())
		binary.Write(&binario7, binary.BigEndian, char)
		writeNextBytes(file, binario7.Bytes())

		file.Seek(superbloque.S_inode_start, 0)
		var binario8 bytes.Buffer
		binary.Write(&binario8, binary.BigEndian, inode)
		writeNextBytes(file, binario8.Bytes())
		var binario9 bytes.Buffer
		binary.Write(&binario9, binary.BigEndian, inodetmp)
		writeNextBytes(file, binario9.Bytes())

		file.Seek(superbloque.S_block_start, 0)
		var binario10 bytes.Buffer
		binary.Write(&binario10, binary.BigEndian, fb)
		writeNextBytes(file, binario10.Bytes())
		var binario11 bytes.Buffer
		binary.Write(&binario11, binary.BigEndian, fileb)
		writeNextBytes(file, binario11.Bytes())
		file.Close()
		fmt.Print("Particion formateada correctamente")
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

package mount

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"unsafe"
)

var path = ""
var name = ""
var disco_amontar = structs.DiscoMontado{}

var path_com = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var name_com = regexp.MustCompile("(?i)\\s?-\\s?name\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")

var ruta_normal = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var ruta_sinext = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var name_id = regexp.MustCompile("([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")

var masterBoot = structs.Mbr{}
var Ebr = structs.Ebr{}

//variables para verificar la existencia del archivo
var pos2 = 0
var abs_path = ""

var alfabeto = [26]string{"A", "B", "C", "D", "E", "F", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func Analizador(input string) {
	if name_com.MatchString(input) && path_com.MatchString(input) {
		path = ruta_normal.FindString(path_com.FindString(input))
		name = name_com.FindString(name_com.FindString(input))
		name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "name", "")
		name = strings.ReplaceAll(name, "=", "")
		name = strings.ReplaceAll(name, " ", "")
		input = path_com.ReplaceAllLiteralString(input, "")
		input = name_com.ReplaceAllLiteralString(input, "")
		fmt.Printf("path = %s\n", path)
		fmt.Printf("name = %s\n", name)
	} else {
		fmt.Println("error sintaxis no esperada los siguientes parametros son obligatorios: ")
		fmt.Println("valores no reconocidos -path: ")
		fmt.Printf("%q\n", path_com.Split(input, -1))
		fmt.Println("valores no reconocidos -name: ")
		fmt.Printf("%q\n", name_com.Split(input, -1))
	}

}

func get_idDislk() structs.DiscoMontado {
	var disco_nuevo = structs.DiscoMontado{}
	disco_nuevo.ID = 1

	if structs.Discos_montados().Len() > 0 {
		for k := structs.Discos_montados().Front(); k != nil; k = k.Next() {
			disco_iterado := structs.DiscoMontado(k.Value.(structs.DiscoMontado))
			var path_guardada = string(disco_iterado.Path[:])
			if strings.Compare(path_guardada, path) == 1 {
				structs.Discos_montados().Remove(k)
				return disco_iterado
			}
			var id = disco_iterado.ID + 1
			disco_nuevo.ID = id
		}

		copy(disco_nuevo.Path[:], path)
		return disco_nuevo
	}
	copy(disco_nuevo.Path[:], path)
	return disco_nuevo
}

func Resultado_Mount() {

	if structs.Discos_montados().Len() > 0 {
		for k := structs.Discos_montados().Front(); k != nil; k = k.Next() {
			disco_iterado := structs.DiscoMontado(k.Value.(structs.DiscoMontado))
			var path_guardada = string(disco_iterado.Path[:])
			fmt.Printf("Disco montado: %s \n", path_guardada)
			fmt.Printf("id Disco montado : %d \n", disco_iterado.ID)
			for s := range disco_iterado.Lista {
				if disco_iterado.Lista[s].EstadoMount != 0 {
					fmt.Printf(" 	particion montada:  %s \n", string(disco_iterado.Lista[s].ID[:]))
				}
			}
		}
	}
	var disco_montado = get_idDislk()
	Abrir_mbr()
	var particion_externa = -1
	var particion_montada = structs.ParticionMontada{}

	for k := range masterBoot.Tabla {
		var name_part = string(masterBoot.Tabla[k].Name[:])
		if string(masterBoot.Tabla[k].Type) == "e" {

			if strings.Compare(name_part, name) == 1 {
				fmt.Println("Error no se puede montar una particion externa")
				return
			}
		}

		if strings.Compare(name_part, name) == 1 {
			particion_montada.EstadoMount = 1
			copy(particion_montada.Nombre[:], name)
			for l := range disco_montado.Lista {
				if disco_montado.Lista[l].EstadoMount == 0 {
					id_disk := strconv.Itoa(disco_montado.ID)
					var id = "23" + id_disk + alfabeto[l]
					copy(particion_montada.ID[:], id)
					disco_montado.Lista[l] = particion_montada
					structs.Montar_disco(disco_montado)
					fmt.Printf("particion montado exitosamente con id %s \n", id)
					return
				}
			}

		}
	}
	for s := range masterBoot.Tabla {
		if string(masterBoot.Tabla[s].Type) == "e" {
			var name_part = string(masterBoot.Tabla[s].Name[:])
			if strings.Compare(name_part, name) == 1 {
				particion_externa = s
				break
			}
		}
	}

	if particion_externa != -1 {
		var inicio = masterBoot.Tabla[particion_externa].Start
		Abrir_ebr(inicio)
		if Ebr.Status == 1 {
			var name_ebr = ""
			for Ebr.Next != -1 {
				Abrir_ebr(Ebr.Next)
				name_ebr = string(Ebr.Name[:])
				if strings.Compare(name_ebr, name) == 1 {
					particion_montada.EstadoMount = 1
					copy(particion_montada.Nombre[:], name)
					for l := range disco_montado.Lista {
						if disco_montado.Lista[l].EstadoMount == 0 {
							id_disk := strconv.Itoa(disco_montado.ID)
							var id = "23" + id_disk + alfabeto[l]
							copy(particion_montada.ID[:], id)
							disco_montado.Lista[l] = particion_montada
							structs.Montar_disco(disco_montado)
							fmt.Printf("particion montada exitosamente con id %s \n", id)
							return
						}
					}
				}
			}
			name_ebr = string(Ebr.Name[:])
			if strings.Compare(name_ebr, name) == 1 {
				particion_montada.EstadoMount = 1
				copy(particion_montada.Nombre[:], name)
				for l := range disco_montado.Lista {
					if disco_montado.Lista[l].EstadoMount == 0 {
						id_disk := strconv.Itoa(disco_montado.ID)
						var id = "23" + id_disk + alfabeto[l]
						copy(particion_montada.ID[:], id)
						disco_montado.Lista[l] = particion_montada
						structs.Montar_disco(disco_montado)
						fmt.Printf("particion montada exitosamente con id %s \n", id)
						return
					}
				}
			}
		}
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

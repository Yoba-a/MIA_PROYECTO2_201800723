package rep

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"unsafe"
)

var path = ""
var name string = ""
var id = ""
var disco_amontar = structs.DiscoMontado{}

var path_com = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.jpg|.png|.gif)")
var name_com = regexp.MustCompile("(?i)\\s?-\\s?name\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var id_com = regexp.MustCompile("(?i)\\s?-\\s?id\\s?=\\s?([0-9]{3}[A-Z])")

var ruta_normal = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.jpg|.png|.gif)")
var ruta_dsk = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var ruta_sinext = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var name_rep = regexp.MustCompile("([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var name_id = regexp.MustCompile("([0-9]{3}[A-Z])")

var masterBoot = structs.Mbr{}
var Ebr = structs.Ebr{}

//variables para verificar la existencia del archivo
var pos2 = 0
var abs_path = ""

var path_disco = ""

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
		id = name_id.FindString(id_com.FindString(input))
		fmt.Printf("path = %s\n", path)
		fmt.Printf("name = %s\n", name)
		fmt.Printf("id = %s\n", id)
	} else {
		fmt.Println("error sintaxis no esperada los siguientes parametros son obligatorios: ")
		fmt.Println("valores no reconocidos -path: ")
		fmt.Printf("%q\n", path_com.Split(input, -1))
		fmt.Println("valores no reconocidos -name: ")
		fmt.Printf("%q\n", name_com.Split(input, -1))
	}

}

func Reportes() {
	if name == "mbr" {
		Reporte_mbr()
	}

}

func Reporte_mbr() {
	var dato string = string(id[2])
	disk, _ := strconv.Atoi(dato)
	if structs.Discos_montados().Len() > 0 {
		var disco_nuevo = structs.DiscoMontado{}

		for k := structs.Discos_montados().Front(); k != nil; k = k.Next() {
			disco_iterado := structs.DiscoMontado(k.Value.(structs.DiscoMontado))

			if disco_iterado.ID == disk {
				disco_nuevo = disco_iterado

				break
			}
		}
		var name_partition = ""
		for k := range disco_nuevo.Lista {
			var id_part_mount = string(disco_nuevo.Lista[k].ID[:])
			if id_part_mount == id {
				name_partition = string(disco_nuevo.Lista[k].Nombre[:])
				path_disco = string(disco_nuevo.Path[:])
				path_disco = ruta_dsk.FindString(path_disco)
				fmt.Println(path_disco)
				break
			}
		}
		if name_partition != "" {
			for pos, char := range path {
				if char == '/' {
					pos2 = pos
				}
			}
			abs_path = path[:pos2]
			fmt.Print(ArchivoExiste(abs_path))
			//creacion del archivo
			if !ArchivoExiste(abs_path) {
				var err = os.Mkdir(abs_path, 0755)
				if err != nil {
					// Aquí puedes manejar mejor el error, es un ejemplo
					panic(err)
				}
			}
			var path_dot = abs_path + "/" + "mbr" + ".dot"
			Abrir_mbr()
			var part1 = name_rep.FindString(string(masterBoot.Tabla[0].Name[:]))
			var part2 = name_rep.FindString(string(masterBoot.Tabla[1].Name[:]))
			var part3 = name_rep.FindString(string(masterBoot.Tabla[2].Name[:]))
			var part4 = name_rep.FindString(string(masterBoot.Tabla[3].Name[:]))
			var tipo byte = 'p'
			var code_dot = "digraph G { \n" +
				"ordering = out \n" +
				"forcelabels=true \n" +
				"graph[ranksep=1,margin=0.3  ]; \n" +
				"node [shape = plaintext];\n " +
				"1 [ label = <<TABLE color = \"black\"> \n" +
				"<TR>\n" +
				"<td > mbr tamaño_disco= " + strconv.Itoa(int(masterBoot.Tamano)) + "</td>\n"
			//por par![](../../../proyecto_archivos/231A.jpg)ticion
			if masterBoot.Tabla[0].Type == tipo {
				var porcentaje = strconv.Itoa(int(masterBoot.Tabla[0].Size) * 100 / int(masterBoot.Tamano))
				code_dot += "<td >" + part1 + "\n " + porcentaje + "%" + "</td>\n"
			} else {
				tam, code := infologicas()
				var colspan = strconv.Itoa(tam)
				code_dot += "<td coslspan=\"" + colspan + "\"" + ">" + "extendida" + "\n " + "</td>\n"
				code_dot += code

			}
			if masterBoot.Tabla[1].Type == tipo {
				var porcentaje = strconv.Itoa(int(masterBoot.Tabla[1].Size) * 100 / int(masterBoot.Tamano))
				code_dot += "<td >" + part2 + "\n " + porcentaje + "%" + "</td>\n"
			} else {
				tam, code := infologicas()
				var colspan = strconv.Itoa(tam)
				code_dot += "<td coslspan=\"" + colspan + "\"" + ">" + "extendida" + "\n " + "</td>\n"
				code_dot += code
			}
			if masterBoot.Tabla[2].Type == tipo {
				var porcentaje = strconv.Itoa(int(masterBoot.Tabla[2].Size) * 100 / int(masterBoot.Tamano))
				code_dot += "<td >" + part3 + "\n " + porcentaje + "%" + "</td>\n"

			} else {
				tam, code := infologicas()
				var colspan = strconv.Itoa(tam)
				code_dot += "<td coslspan=\"" + colspan + "\"" + ">" + "extendida" + "\n " + "</td>\n"
				code_dot += code
			}
			if masterBoot.Tabla[3].Type == tipo {
				var porcentaje = strconv.Itoa(int(masterBoot.Tabla[3].Size) * 100 / int(masterBoot.Tamano))
				code_dot += "<td >" + part4 + "\n " + porcentaje + "%" + "</td>\n"

			} else {
				tam, code := infologicas()
				var colspan = strconv.Itoa(tam)
				code_dot += "<td coslspan=\"" + colspan + "\"" + ">" + "extendida" + "\n " + "</td>\n"
				code_dot += code

			}

			code_dot += "</TR>\n" +
				"</TABLE>> dir =none color=white style =none]\n" +
				"}"
			f, err := os.Create(path_dot)
			check(err)
			defer f.Close()

			f.Sync()
			w := bufio.NewWriter(f)
			n4, err := w.WriteString(code_dot)
			check(err)
			fmt.Printf("escribi estos %d bytes \n", n4)
			w.Flush()
			dot := "dot"
			format := "-Tjpg"
			dot_file := path_dot
			ouput := "-o"
			ab_pa := abs_path + "/" + id + ".jpg"
			cmd := exec.Command(dot, format, dot_file, ouput, ab_pa)

			stdout, err := cmd.Output()

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// Print the output
			fmt.Println(string(stdout))
		} else {
			fmt.Println("No hay ni una particion montada con ese id")
		}
	} else {
		fmt.Println("No hay ni una particion montada")
	}
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func Abrir_mbr() {
	for pos, char := range path_disco {
		if char == '/' {
			pos2 = pos
		}
	}
	abs_path = path_disco[:pos2]
	fmt.Print(abs_path)
	if ArchivoExiste(path_disco) {
		file, err := os.Open(path_disco)
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

func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("ERROR  A LEER BYTES", err)
	}

	return bytes
}

func Abrir_ebr(inicio_ebr int64) {
	for pos, char := range path_disco {
		if char == '/' {
			pos2 = pos
		}
	}
	abs_path = path_disco[:pos2]

	if ArchivoExiste(path_disco) {
		file, err := os.Open(path_disco)
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

func infologicas() (int, string) {
	var size int = 0
	var info = ""
	var particion_externa int64 = -1
	for s := range masterBoot.Tabla {
		if string(masterBoot.Tabla[s].Type) == "e" {
			particion_externa = masterBoot.Tabla[s].Start
			break

		}
	}

	Abrir_ebr(particion_externa)

	for Ebr.Next != -1 {
		Abrir_ebr(Ebr.Next)
		var name_ebr = name_rep.FindString(string(Ebr.Name[:]))
		var porcentaje = strconv.Itoa(int(Ebr.Size) * 100 / int(masterBoot.Tamano))
		info += "<td >ebr</td>\n"
		info += "<td >" + name_ebr + "\n " + porcentaje + "%" + "</td>\n"
		size++
	}
	size++
	var name_ebr = name_rep.FindString(string(Ebr.Name[:]))
	var porcentaje = strconv.Itoa(int(Ebr.Size) * 100 / int(masterBoot.Tamano))
	info += "<td >ebr</td>\n"
	info += "<td >" + name_ebr + "\n " + porcentaje + "%" + "</td>\n"

	return size, info
}

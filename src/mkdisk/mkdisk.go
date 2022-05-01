package mkdisk

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"time"
)

//struct usado para el mbr

//variables globales utilizadas para resguardar los parametros de entrada
var Size = ""
var fit = ""
var unit = ""
var path = ""

//variables utilizadas para analizar las entradas
var path_com = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var size_com = regexp.MustCompile("(?i)\\s?-\\s?size\\s?=\\s?[0-9]+")
var fit_com = regexp.MustCompile("(?i)\\s?-\\s?fit\\s?=\\s?(bf|ff|wf)")
var unit_com = regexp.MustCompile("(?i)\\s?-\\s?unit\\s?=\\s?(k|m)")
var numeros = regexp.MustCompile("[0-9]+")
var ruta_normal = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var ruta_sinext = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var fit_values = regexp.MustCompile("(bf|ff|wf)")
var unit_val = regexp.MustCompile("(k|m)")
var barra = regexp.MustCompile("/")

//analizador nos sirve para obtener los parametros de la entrada
func Analizador2(input string) {
	if path_com.MatchString(input) {
		path = ruta_normal.FindString(path_com.FindString(input))
		fmt.Printf("path = %s\n", path)
	} else {
		fmt.Println("error sintaxis no esperada")
		fmt.Println("valores no reconocidos -path: ")
		fmt.Printf("%q\n", path_com.Split(input, -1))
	}
}
func Analizador(input string) {
	if size_com.MatchString(input) && path_com.MatchString(input) {
		Size = numeros.FindString(size_com.FindString(input))
		path = ruta_normal.FindString(path_com.FindString(input))
		input = size_com.ReplaceAllLiteralString(input, "")
		input = path_com.ReplaceAllLiteralString(input, "")
		fmt.Printf("size = %s\n", Size)
		fmt.Printf("path = %s\n", path)
		if regexp.MustCompile("(?i)fit").MatchString(input) {
			if fit_com.MatchString(input) {
				fit = fit_values.FindString(fit_com.FindString(input))
				fmt.Printf("fit = %s\n", fit)
			} else {
				fmt.Println("error sintaxis no esperada")
				fmt.Println("valores no reconocidos -fit: ")
				fmt.Printf("%q\n", fit_com.Split(input, -1))
			}
		} else {
			fit = "wf"
			fmt.Printf("fit = %s\n", fit)
		}
		if regexp.MustCompile("(?i)unit").MatchString(input) {
			if unit_com.MatchString(input) {
				unit = unit_val.FindString(unit_com.FindString(input))
				fmt.Printf("unit = %s\n", unit)
			} else {
				fmt.Println("error sintaxis no esperada")
				fmt.Println("valores no reconocidos -unit: ")
				fmt.Printf("%q\n", unit_com.Split(input, -1))
			}
		} else {
			unit = "m"
			fmt.Printf("unit = %s\n", unit)
		}
	} else {
		fmt.Println("error sintaxis no esperada")
		fmt.Println("valores no reconocidos -path: ")
		fmt.Printf("%q\n", path_com.Split(input, -1))
		fmt.Println("valores no reconocidos -size: ")
		fmt.Printf("%q\n", size_com.Split(input, -1))

	}
}

var pos2 = 0
var abs_path = ""

//esta funcion nos servira para crear un disco con su mbr
func CrearDisco() {

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
	archivo, err := os.Create(path)

	size, err := strconv.ParseInt(Size, 10, 64)
	fmt.Println(size, err, reflect.TypeOf(size))
	defer archivo.Close()
	if err != nil {
		log.Fatal(err)
	}
	var vacio int8 = 0
	s := &vacio
	var num int64 = 0
	//Definiendo tamano
	if strings.Compare(strings.ToLower(unit), "m") == 0 {
		num = int64(size) * 1024 * 1024
	} else if strings.Compare(strings.ToLower(unit), "k") == 0 {
		num = int64(size) * 1024
	}
	num = num - 1
	//Llenando el archivo

	//colocando el primer byte
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	writeNextBytes(archivo, binario.Bytes())

	//situando el cursor en la ultima posicion
	archivo.Seek(num, 0)

	//colocando el ultimo byte para rellenar
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	writeNextBytes(archivo, binario2.Bytes())

	//Regresando el cursor a 0 para escribir el mbr
	archivo.Seek(0, 0)

	//Formando el MBR
	disco := structs.Mbr{}
	disco.Tamano = num + 1

	fechahora := time.Now()
	fechahoraArreglo := strings.Split(fechahora.String(), "")
	fechahoraCadena := ""
	for i := 0; i < 16; i++ {
		fechahoraCadena = fechahoraCadena + fechahoraArreglo[i]
	}
	copy(disco.Fecha[:], fechahoraCadena)
	copy(disco.Fit[:], fit)
	var signature int8
	binary.Read(rand.Reader, binary.LittleEndian, &signature)
	if signature < 0 {
		signature = signature * -1
	}
	disco.Firma = signature

	//Escribiendo el MBR
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, disco)
	writeNextBytes(archivo, binario3.Bytes())
	//path := path
	//graficarDISCO(path)
	//graficarMBR(path)
}

func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func EliminarDisco() {
	if ArchivoExiste(path) {
		err := os.Remove(path)
		if err != nil {
			fmt.Printf("Error eliminando disco: %v\n", err)
		} else {
			fmt.Println("Disco eliminado correctamente")
		}
	}
}

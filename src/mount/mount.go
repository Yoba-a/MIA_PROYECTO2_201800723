package mount

import (
	"fmt"
	"regexp"
	"strings"
	"structs"
)

var path = ""
var name = ""
var disco_amontar = structs.DiscoMontado{}

var path_com = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var name_com = regexp.MustCompile("(?i)\\s?-\\s?name\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")

var ruta_normal = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var ruta_sinext = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var name_id = regexp.MustCompile("([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")

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
	var ultimo_id = 1

	if structs.Discos_montados().Len() > 0 {
		for k := structs.Discos_montados().Front(); k != nil; k = k.Next() {
			disco_iterado := structs.DiscoMontado(k.Value.(structs.DiscoMontado))
			var path_guardada = string(disco_iterado.Path[:])
			if strings.Compare(path_guardada, path) == 1 {
				return disco_iterado
			}
			ultimo_id = disco_iterado.ID + 1
		}
		disco_nuevo.ID = ultimo_id
		copy(disco_nuevo.Path[:], path)
		return disco_nuevo
	}
	return disco_nuevo
}

func resultado_Mount() {

}

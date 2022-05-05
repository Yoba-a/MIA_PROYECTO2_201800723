package mkfs

import (
	"fmt"
	"regexp"
	"structs"
)

var path = ""
var name string = ""
var id = ""
var disco_amontar = structs.DiscoMontado{}
var tipo_ = ""

var type_com = regexp.MustCompile("(?i)\\s?-\\s?type\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var id_com = regexp.MustCompile("(?i)\\s?-\\s?id\\s?=\\s?([0-9]{3}[A-Z])")

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
		fmt.Printf("path = %s\n", path)
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

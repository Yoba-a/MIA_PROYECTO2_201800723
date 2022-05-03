package fdisk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"unsafe"
)

var size = ""
var fit = ""
var unit = ""
var path = ""
var type_ = ""
var name = ""

///PARA OBTENER LOS VALORES DE LA ENTRADA
var size_com = regexp.MustCompile("(?i)\\s?-\\s?size\\s?=\\s?[0-9]+")
var fit_com = regexp.MustCompile("(?i)\\s?-\\s?fit\\s?=\\s?(bf|ff|wf)")
var unit_com = regexp.MustCompile("(?i)\\s?-\\s?unit\\s?=\\s?(k|m)")
var path_com = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var type_com = regexp.MustCompile("(?i)\\s?-\\s?type\\s?=\\s?(P|p|E|e|L|l)")
var type_aux = regexp.MustCompile("(?i)\\s?-\\s?type\\s?=\\s?")
var name_com = regexp.MustCompile("(?i)\\s?-\\s?name\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")

///AUXILIARES
var numeros = regexp.MustCompile("[0-9]+")
var ruta_normal = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var ruta_sinext = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var fit_values = regexp.MustCompile("(bf|ff|wf)")
var unit_val = regexp.MustCompile("(k|m)")
var type_val = regexp.MustCompile("(P|p|E|e|L|l)")
var name_id = regexp.MustCompile("([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var barra = regexp.MustCompile("/")

var masterBoot = structs.Mbr{}
var Ebr = structs.Ebr{}

func Analizador(input string) {
	if size_com.MatchString(input) && path_com.MatchString(input) && name_com.MatchString(input) {
		size = numeros.FindString(size_com.FindString(input))
		path = ruta_normal.FindString(path_com.FindString(input))
		name = name_com.FindString(name_com.FindString(input))
		name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "name", "")
		name = strings.ReplaceAll(name, "=", "")
		name = strings.ReplaceAll(name, " ", "")
		input = size_com.ReplaceAllLiteralString(input, "")
		input = path_com.ReplaceAllLiteralString(input, "")
		input = name_com.ReplaceAllLiteralString(input, "")
		fmt.Printf("size = %s\n", size)
		fmt.Printf("path = %s\n", path)
		fmt.Printf("name = %s\n", name)
		if regexp.MustCompile("(?i)fit").MatchString(input) {
			if fit_com.MatchString(input) {
				fit = fit_values.FindString(fit_com.FindString(input))
				input = fit_com.ReplaceAllLiteralString(input, "")
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
				input = unit_com.ReplaceAllLiteralString(input, "")
				fmt.Printf("unit = %s\n", unit)
			} else {
				fmt.Println("error sintaxis no esperada")
				fmt.Println("valores no reconocidos -unit: ")
				fmt.Printf("%q\n", unit_com.Split(input, -1))
			}
		} else {
			unit = "K"
			fmt.Printf("unit = %s\n", unit)
		}
		if regexp.MustCompile("(?i)type").MatchString(input) {
			if type_com.MatchString(input) {
				input = type_aux.ReplaceAllLiteralString(input, "")
				type_ = strings.ToLower(type_val.FindString(input))
				fmt.Printf("type = %s\n", type_)
			} else {
				fmt.Println("error sintaxis no esperada")
				fmt.Println("valores no reconocidos -type: ")
				fmt.Printf("%q\n", type_com.Split(input, -1))
			}
		} else {
			type_ = "p"
			fmt.Printf("type = %s\n", type_)
		}

	} else {
		fmt.Println("error sintaxis no esperada los siguientes parametros son obligatorios: ")
		fmt.Println("valores no reconocidos -path: ")
		fmt.Printf("%q\n", path_com.Split(input, -1))
		fmt.Println("valores no reconocidos -size: ")
		fmt.Printf("%q\n", size_com.Split(input, -1))
		fmt.Println("valores no reconocidos -name: ")
		fmt.Printf("%q\n", name_com.Split(input, -1))
	}

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

//variables para verificar la existencia del archivo
var pos2 = 0
var abs_path = ""

func determinar_particion(size_newP int64) [3]int64 {
	data := [3]int64{-1, -1, -1}
	var p1_fin = masterBoot.Tabla[0].Start + masterBoot.Tabla[0].Size
	var p2_fin = masterBoot.Tabla[1].Start + masterBoot.Tabla[1].Size
	var p3_fin = masterBoot.Tabla[2].Start + masterBoot.Tabla[2].Size

	if masterBoot.Tabla[0].Size == 0 || masterBoot.Tabla[1].Size == 0 ||
		masterBoot.Tabla[2].Size == 0 || masterBoot.Tabla[3].Size == 0 {
		if masterBoot.Tabla[0].Size == 0 && masterBoot.Tabla[0].Status == 0 {
			data[0] = 0
			data[1] = size_newP + int64(unsafe.Sizeof(structs.Mbr{}))
			data[2] = int64(unsafe.Sizeof(structs.Mbr{})) + 1
		} else if masterBoot.Tabla[1].Size == 0 && masterBoot.Tabla[1].Status == 0 {
			data[0] = 1
			data[1] = masterBoot.Tabla[0].Size + size_newP + int64(unsafe.Sizeof(structs.Mbr{}))
			data[2] = p1_fin + 1
		} else if masterBoot.Tabla[2].Size == 0 && masterBoot.Tabla[2].Status == 0 {
			data[0] = 2
			data[1] = masterBoot.Tabla[0].Size + masterBoot.Tabla[1].Size + size_newP + int64(unsafe.Sizeof(structs.Mbr{}))
			data[2] = p2_fin + 1
		} else if masterBoot.Tabla[3].Size == 0 && masterBoot.Tabla[3].Status == 0 {
			data[0] = 3
			data[1] = masterBoot.Tabla[0].Size + masterBoot.Tabla[1].Size + masterBoot.Tabla[2].Size + size_newP + int64(unsafe.Sizeof(structs.Mbr{}))
			data[2] = p3_fin + 1
		} else {
			fmt.Println("No hay espacio para tu particinon en este disco")
		}
	}

	return data
}

func Abrir_mbr() {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	abs_path = path[:pos2]
	fmt.Print(abs_path)
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
		fmt.Println("Mbr Tamano:", masterBoot.Tamano)
		fmt.Println("Mbr Fecha creacion:", string(masterBoot.Fecha[:]))
		fmt.Println("Mbr Disk Signarue:", masterBoot.Firma)
		fmt.Println("Disco Fit:", string(masterBoot.Fit[:]))
		for k := range masterBoot.Tabla {
			fmt.Println("particion:", string(masterBoot.Tabla[k].Name[:]))
			fmt.Println("size :", masterBoot.Tabla[k].Size)
		}
		file.Close()
		crear_particion()
	} else {
		fmt.Print("error el disco no existe en esta computadora...")
	}
}

func get_size() int64 {
	var num int64 = 0
	Size, err := strconv.ParseInt(size, 10, 64)
	fmt.Println(size, err, reflect.TypeOf(size))
	if err != nil {
		log.Fatal(err)
	}
	//Definiendo tamano
	if strings.Compare(strings.ToLower(unit), "m") == 0 {
		num = int64(Size) * 1024 * 1024
	} else if strings.Compare(strings.ToLower(unit), "k") == 0 {
		num = int64(Size) * 1024
	}
	num = num - 1
	return num
}

func nombre_existente() bool {
	var existe = false
	for k := range masterBoot.Tabla {
		name_aux := string(masterBoot.Tabla[k].Name[:])
		fmt.Println(name_aux)
		if strings.Compare(name_aux, name) == 1 {
			existe = true
			break
		}
	}
	return existe
}

func verificar_p_ext() bool {
	var ext bool = false
	for k := range masterBoot.Tabla {
		tipo_ext := string(masterBoot.Tabla[k].Type)
		if tipo_ext == "e" {
			ext = true
		}
	}
	return ext
}

func valores_particionExt() [3]int64 {
	data := [3]int64{-1, -1, -1}
	for k := range masterBoot.Tabla {
		if string(masterBoot.Tabla[k].Type) == "e" {
			data[0] = int64(k)
			data[1] = masterBoot.Tabla[k].Start
			data[2] = masterBoot.Tabla[k].Size
			break
		}
	}

	return data
}

func crear_particion() {
	if particiones_vacias() {
		if type_ != "l" {
			masterBoot.Tabla[0].Size = get_size()
			if masterBoot.Tabla[0].Size > 0 && masterBoot.Tabla[0].Size < masterBoot.Tamano {
				masterBoot.Tabla[0].Status = 1
				copy(masterBoot.Tabla[0].Fit[:], fit)
				if type_ == "p" {
					masterBoot.Tabla[0].Type = 'p'
				} else if type_ == "e" {
					masterBoot.Tabla[0].Type = 'e'
					var size_mbr = int64(unsafe.Sizeof(masterBoot))
					insertar_ebr(size_mbr)
				}
				masterBoot.Tabla[0].Start = int64(unsafe.Sizeof(structs.Mbr{})) + 1
				copy(masterBoot.Tabla[0].Name[:], name)
				insertar_mbr()
			} else {
				fmt.Println("(particion) nose puede crear una particion ingrese un tamanio correcto")
			}
		} else {
			fmt.Println("(particion) nose puede crear una particion logica sin haber creado una extendida antes")
		}
	} else {
		if type_ != "l" {
			var size_parti int64 = get_size()
			if size_parti < masterBoot.Tamano {
				var resultado [3]int64 = determinar_particion(size_parti)
				var particion_aModificar = resultado[0]
				var tamano_usado = resultado[1]
				var inicio_particion = resultado[2]
				if particion_aModificar != -1 && tamano_usado != -1 && inicio_particion != -1 {
					if !nombre_existente() {
						if tamano_usado < masterBoot.Tamano {
							if type_ == "p" {
								masterBoot.Tabla[particion_aModificar].Status = 1
								copy(masterBoot.Tabla[particion_aModificar].Fit[:], fit)
								masterBoot.Tabla[particion_aModificar].Type = 'p'
								masterBoot.Tabla[particion_aModificar].Start = inicio_particion
								masterBoot.Tabla[particion_aModificar].Size = get_size()
								copy(masterBoot.Tabla[particion_aModificar].Name[:], name)
								insertar_mbr()
								fmt.Println("(partition) particion primaria creada exitosamente")
							} else if !verificar_p_ext() && type_ == "e" {
								masterBoot.Tabla[particion_aModificar].Status = 1
								copy(masterBoot.Tabla[particion_aModificar].Fit[:], fit)
								masterBoot.Tabla[particion_aModificar].Type = 'e'
								masterBoot.Tabla[particion_aModificar].Start = inicio_particion
								masterBoot.Tabla[particion_aModificar].Size = get_size()
								copy(masterBoot.Tabla[particion_aModificar].Name[:], name)
								insertar_mbr()
								insertar_ebr(inicio_particion + 1)
								fmt.Println("(partition) particion extendida creada exitosamente")
							} else {
								fmt.Println("(fdisk_error) Solo se puede tener una particion extendida ")
							}
						} else {
							fmt.Println("(fdisk_error) el tamaño deseado excede el tamaño del disco ")
						}
					} else {
						fmt.Println("(partition) El nombre de la particion ya esta ocupado escriba otro ")
					}
				} else {
					fmt.Println("(partition) no se pudo crear la particion por que no hay particiones vacias")
				}

			}
		} else {
			fmt.Println("Creando particion logica")
			if verificar_p_ext() {
				var result = valores_particionExt()
				var indice = result[0]
				var inicio = result[1]
				var tamano_ext = result[2]
				if indice != -1 && inicio != -1 && tamano_ext != -1 && masterBoot.Tabla[indice].Status != 0 {
					Abrir_ebr(inicio)
					if Ebr.Status == 1 {
						var name_repetido = false
						var name_ebr = ""
						for Ebr.Next != -1 {
							Abrir_ebr(Ebr.Next)
							name_ebr = string(Ebr.Name[:])
							if strings.Compare(name_ebr, name) == 1 {
								name_repetido = true
								break
							}
						}
						if !name_repetido {
							name_ebr = string(Ebr.Name[:])
							if Ebr.Next == -1 && strings.Compare(name_ebr, name) != 1 {
								Ebr.Next = Ebr.Size + Ebr.Start + 1
								var inicio_nueva = Ebr.Next
								if (inicio_nueva + get_size() + 1) < tamano_ext {
									insertar_ebr(Ebr.Start)
									copy(Ebr.Fit[:], fit)
									copy(Ebr.Name[:], name)
									Ebr.Start = inicio_nueva
									Ebr.Size = get_size()
									Ebr.Next = -1
									Ebr.Status = 1
									insertar_ebr(inicio_nueva)
								} else {
									fmt.Println("error no hay espacio suficiente en la particion extendida para esta particion logica ")
								}
							} else {
								fmt.Println("no se puede crear la particion logica debido a que hay otra con el mismo nombre")
							}
						} else {
							fmt.Println("no se puede crear la particion logica debido a que hay otra con el mismo nombre")
						}
					} else {
						copy(Ebr.Fit[:], fit)
						copy(Ebr.Name[:], name)
						Ebr.Start = inicio
						Ebr.Size = get_size()
						Ebr.Next = -1
						Ebr.Status = 1
						insertar_ebr(inicio)
					}
				} else {
					fmt.Println("error no existe una particion extendida, cree una ")
				}
			} else {
				fmt.Println("error no existe una particion extendida, cree una ")
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

func particiones_vacias() bool {
	var resultado = true
	for _, s := range masterBoot.Tabla {
		if s.Status == 1 {
			resultado = false
			break
		}
	}
	return resultado
}

func insertar_mbr() {
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	mbr := masterBoot

	file.Seek(0, 0)
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, mbr)
	writeNextBytes(file, binario3.Bytes())
	file.Close()
}

func insertar_ebr(size int64) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(size, 0)
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, Ebr)
	writeNextBytes(file, binario3.Bytes())
	file.Close()
	fmt.Print("Ebr insertado correctamente ")
}

func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("ERROR  A LEER BYTES", err)
	}

	return bytes
}

package main

import (
	"bufio"
	"fdisk"
	"fmt"
	"log"
	"mkdisk"
	mkfs2 "mkfs"
	"mount"
	"os"
	"regexp"
	"rep"
	"strings"
)

var mkdisk_com = regexp.MustCompile("(?i)mkdisk")
var rmdisk_com = regexp.MustCompile("(?i)rmdisk")
var read = regexp.MustCompile("(?i)read")
var fdisk_com = regexp.MustCompile("(?i)fdisk")
var mount_com = regexp.MustCompile("(?i)mount")
var exec_com = regexp.MustCompile("(?i)exec")
var rep_com = regexp.MustCompile("(?i)rep")
var mkfs = regexp.MustCompile("(?i)mkfs")
var comentario = regexp.MustCompile("(?i)#")

var path_exec = ""

//para leer el path en exec
var path_com = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.sh|.txt|.script)")
var ruta_normal = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.sh|.txt|.script)")

func main() {

	menu :=
		`
	------------------------------INGRESE UN COMANDO------------------------------
	--------------------------------exit para salir-------------------------------
	
>`

	reader := bufio.NewReader(os.Stdin)

	// Leer hasta el separador de salto de línea

	for {
		fmt.Print(menu)
		entrada, _ := reader.ReadString('\n')
		eleccion := strings.TrimRight(entrada, "\n") // Remover el salto de línea de la entrada del usuario

		if eleccion == "exit" {
			fmt.Println("Adios!")
			break
		} else if exec_com.MatchString(eleccion) {
			analizador_exec(eleccion)
			if ArchivoExiste(path_exec) {
				file, err := os.Open(path_exec)
				//handle errors while opening
				if err != nil {
					log.Fatalf("Error when opening file: %s", err)
				}
				fileScanner := bufio.NewScanner(file)
				// read line by line
				for fileScanner.Scan() {
					if fileScanner.Text() == "pause" {
						fmt.Println("pausa")
						reader2 := bufio.NewReader(os.Stdin)
						entrada2, _ := reader2.ReadString('\n')
						eleccion2 := strings.TrimRight(entrada2, "\n") // Remover el salto de línea de la entrada del usuario
						fmt.Println(eleccion2)
						fmt.Println("continuando")
					}
					ele := strings.TrimRight(fileScanner.Text(), "\n")
					logic(ele[0:])
				}
				// handle first encountered error while reading
				if err := fileScanner.Err(); err != nil {
					log.Fatalf("Error while reading file: %s", err)
				}
				file.Close()
			} else {
				fmt.Println("el archivo sh no existe!")
			}

		} else {
			fmt.Printf("%q\n", mkdisk_com.FindString(eleccion))
			logic(eleccion)
		}

	}
}

func analizador_exec(input string) {
	if path_com.MatchString(input) {
		path_exec = ruta_normal.FindString(path_com.FindString(input))
		input = path_com.ReplaceAllLiteralString(input, "")
		fmt.Printf("path = %s\n", path_exec)
	} else {
		fmt.Println("error sintaxis no esperada los siguientes parametros son obligatorios: ")
		fmt.Println("valores no reconocidos -path: ")
		fmt.Printf("%q\n", path_com.Split(input, -1))
	}
}

func logic(eleccion string) {
	fmt.Println(eleccion)
	if comentario.MatchString(eleccion) {
		eleccion = comentario.ReplaceAllLiteralString(eleccion, "")
		fmt.Println(eleccion)
	} else if mkdisk_com.MatchString(eleccion) {
		fmt.Println("comando mkdisk, creacion de disco en proceso...")
		eleccion = mkdisk_com.ReplaceAllLiteralString(eleccion, "")
		mkdisk.Analizador(eleccion)
		mkdisk.CrearDisco()
	} else if rmdisk_com.MatchString(eleccion) {
		fmt.Println("comando rmkdisk, eliminacion de disco en proceso...")
		eleccion = rmdisk_com.ReplaceAllLiteralString(eleccion, "")
		mkdisk.Analizador2(eleccion)
		mkdisk.EliminarDisco()
	} else if fdisk_com.MatchString(eleccion) {
		fmt.Println("comando fdisk, creacion de particion en proceso...")
		eleccion = fdisk_com.ReplaceAllLiteralString(eleccion, "")
		fdisk.Analizador(eleccion)

		fdisk.Abrir_mbr()
		fdisk.Error = false

	} else if mount_com.MatchString(eleccion) {
		fmt.Println("comando mount, montando particion en proceso...")
		eleccion = mount_com.ReplaceAllLiteralString(eleccion, "")
		mount.Analizador(eleccion)
		mount.Resultado_Mount()
	} else if rep_com.MatchString(eleccion) {
		fmt.Println("comando rep, reporte en proceso...")
		eleccion = mount_com.ReplaceAllLiteralString(eleccion, "")
		rep.Analizador(eleccion)
		rep.Reportes()
	} else if mkfs.MatchString(eleccion) {
		fmt.Println("comando mkfs, reporte en proceso...")
		eleccion = mkfs.ReplaceAllLiteralString(eleccion, "")
		mkfs2.Analizador(eleccion)
		mkfs2.FormatearExt2()
	}
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

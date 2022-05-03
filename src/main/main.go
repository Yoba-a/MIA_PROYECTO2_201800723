package main

import (
	"bufio"
	"fdisk"
	"fmt"
	"mkdisk"
	"mount"
	"os"
	"regexp"
	"strings"
)

var mkdisk_com = regexp.MustCompile("(?i)mkdisk")
var rmdisk_com = regexp.MustCompile("(?i)rmdisk")
var read = regexp.MustCompile("(?i)read")
var fdisk_com = regexp.MustCompile("(?i)fdisk")
var mount_com = regexp.MustCompile("(?i)mount")

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
		}

		fmt.Printf("%q\n", mkdisk_com.FindString(eleccion))

		if mkdisk_com.MatchString(eleccion) {
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

		} else if mount_com.MatchString(eleccion) {
			fmt.Println("comando fdisk, creacion de particion en proceso...")
			eleccion = mount_com.ReplaceAllLiteralString(eleccion, "")
			mount.Analizador(eleccion)
		}

	}
}
